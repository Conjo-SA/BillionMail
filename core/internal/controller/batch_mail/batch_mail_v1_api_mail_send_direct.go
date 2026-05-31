package batch_mail

import (
	"billionmail-core/api/batch_mail/v1"
	"billionmail-core/internal/service/public"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ApiMailSendDirect sends an email without a pre-configured template.
// The subject and HTML body are provided inline in the request.
// Rate limits (daily_limit, monthly_limit) configured on the API key are enforced.
func (c *ControllerV1) ApiMailSendDirect(ctx context.Context, req *v1.ApiMailSendDirectReq) (res *v1.ApiMailSendDirectRes, err error) {
	res = &v1.ApiMailSendDirectRes{}
	clientIP := g.RequestFromCtx(ctx).GetClientIp()

	// 1. Authenticate via API key
	apiTemplate, err := getApiTemplateByKey(ctx, req.ApiKey, clientIP)
	if err != nil {
		res.Code = 1001
		res.SetError(gerror.New(public.LangCtx(ctx, err.Error())))
		return res, nil
	}

	// 2. IP whitelist check
	if err = CheckClientIP(ctx, apiTemplate.Id, clientIP); err != nil {
		res.Code = 1002
		res.SetError(gerror.New(public.LangCtx(ctx, err.Error())))
		return res, nil
	}

	// 3. Validate recipient
	if req.To == "" || !strings.Contains(req.To, "@") {
		res.Code = 1003
		res.SetError(gerror.New(public.LangCtx(ctx, "Invalid recipient")))
		return res, nil
	}

	// 4. Enforce rate limits
	if err = checkApiRateLimit(ctx, apiTemplate.Id, apiTemplate.DailyLimit, apiTemplate.MonthlyLimit); err != nil {
		res.Code = 1006
		res.SetError(gerror.New(err.Error()))
		return res, nil
	}

	// 5. Resolve sender address
	addresser := req.From
	if addresser == "" {
		addresser = apiTemplate.Addresser
	}
	fullName := req.FromName
	if fullName == "" {
		fullName = apiTemplate.FullName
	}

	// 6. Build a transient api template copy with resolved sender name
	tpl := *apiTemplate
	tpl.FullName = fullName

	// 7. Enqueue
	if err = recordApiMailLogFull(ctx, &tpl, req.To, addresser, nil, req.Subject, req.Html); err != nil {
		res.Code = 1005
		res.SetError(gerror.New(public.LangCtx(ctx, "Failed to record email log: {}", err.Error())))
		return res, nil
	}

	res.SetSuccess(public.LangCtx(ctx, "Email queued successfully"))
	return res, nil
}

// checkApiRateLimit enforces daily and monthly send limits using Redis counters.
// 0 means unlimited. Returns an error when the limit is exceeded.
func checkApiRateLimit(ctx context.Context, apiId, dailyLimit, monthlyLimit int) error {
	now := time.Now()

	if dailyLimit > 0 {
		dayKey := fmt.Sprintf("bm:api_rate:daily:%d:%s", apiId, now.Format("2006-01-02"))
		count, incrErr := g.Redis().Incr(ctx, dayKey)
		if incrErr != nil {
			return incrErr
		}
		// Set TTL on first write (expire at next midnight)
		if count == 1 {
			tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			_, _ = g.Redis().Expire(ctx, dayKey, int64(tomorrow.Sub(now).Seconds()))
		}
		if int(count) > dailyLimit {
			_, _ = g.Redis().Decr(ctx, dayKey)
			return gerror.Newf("Daily send limit of %d reached", dailyLimit)
		}
	}

	if monthlyLimit > 0 {
		monthKey := fmt.Sprintf("bm:api_rate:monthly:%d:%s", apiId, now.Format("2006-01"))
		count, incrErr := g.Redis().Incr(ctx, monthKey)
		if incrErr != nil {
			// Roll back daily counter on error
			if dailyLimit > 0 {
				dayKey := fmt.Sprintf("bm:api_rate:daily:%d:%s", apiId, now.Format("2006-01-02"))
				_, _ = g.Redis().Decr(ctx, dayKey)
			}
			return incrErr
		}
		if count == 1 {
			nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			_, _ = g.Redis().Expire(ctx, monthKey, int64(nextMonth.Sub(now).Seconds()))
		}
		if int(count) > monthlyLimit {
			_, _ = g.Redis().Decr(ctx, monthKey)
			if dailyLimit > 0 {
				dayKey := fmt.Sprintf("bm:api_rate:daily:%d:%s", apiId, now.Format("2006-01-02"))
				_, _ = g.Redis().Decr(ctx, dayKey)
			}
			return gerror.Newf("Monthly send limit of %d reached", monthlyLimit)
		}
	}

	return nil
}

// getApiRateCounts returns (dailySent, monthlySent) for an API key from Redis.
func getApiRateCounts(ctx context.Context, apiId int) (int, int) {
	now := time.Now()
	dayKey := fmt.Sprintf("bm:api_rate:daily:%d:%s", apiId, now.Format("2006-01-02"))
	monthKey := fmt.Sprintf("bm:api_rate:monthly:%d:%s", apiId, now.Format("2006-01"))

	dailyVal, _ := g.Redis().Get(ctx, dayKey)
	monthlyVal, _ := g.Redis().Get(ctx, monthKey)

	return dailyVal.Int(), monthlyVal.Int()
}
