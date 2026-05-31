# Direct Send API

The Direct Send endpoint allows you to send transactional emails inline — no pre-configured template required. This is ideal for integrations where the email content is generated dynamically (password resets, OTPs, order confirmations, etc.).

## Endpoint

```
POST /api/batch_mail/api/send_direct
```

Authentication is via the `X-API-Key` header. No JWT or session cookie needed.

## Request

### Headers

| Header          | Required | Description                         |
|-----------------|----------|-------------------------------------|
| `X-API-Key`     | ✅        | Your API key from the Send API page |
| `Content-Type`  | ✅        | Must be `application/json`          |

### Body (JSON)

| Field       | Required | Description                                                  |
|-------------|----------|--------------------------------------------------------------|
| `to`        | ✅        | Recipient email address                                      |
| `subject`   | ✅        | Email subject line                                           |
| `html`      | ✅        | Full HTML body of the email                                  |
| `from`      | ❌        | Sender address (overrides the API default)                   |
| `from_name` | ❌        | Sender display name (overrides the API default)              |

## Example

```bash
curl -X POST 'https://your-domain/api/batch_mail/api/send_direct' \
  -H 'X-API-Key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "to": "user@example.com",
    "subject": "Your verification code",
    "html": "<p>Your code is <strong>123456</strong>.</p>"
  }'
```

## Response

```json
{ "code": 0, "message": "Email queued successfully" }
```

Non-zero `code` values indicate errors (e.g. invalid API key, rate limit exceeded).

## Rate Limits

Each API key can have a **daily limit** and/or a **monthly limit** configured. Set to `0` for unlimited.

When a limit is exceeded the API returns:

```json
{ "code": 50, "message": "Daily send limit of 1000 reached" }
```

Current usage counters are visible in the **Send API** page (Usage column shows `sent/limit` per day · per month).

## Direct Send vs. Template Send

| Feature              | Template Send (`/send`) | Direct Send (`/send_direct`) |
|----------------------|------------------------|------------------------------|
| Template required    | ✅ Yes                  | ❌ No                         |
| Contact saved to DB  | ✅ Yes                  | ❌ No                         |
| Unsubscribe tracking | ✅ Yes                  | ❌ No                         |
| Open/click tracking  | Depends on API config  | ❌ No                         |
| Use case             | Marketing, newsletters  | Transactional, OTP, alerts   |
