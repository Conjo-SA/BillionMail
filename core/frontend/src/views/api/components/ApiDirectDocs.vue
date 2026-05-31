<template>
	<n-card :content-style="{ padding: '0' }">
		<n-collapse>
			<n-collapse-item :title="$t('api.directSend.collapseTitle')" name="direct-send">
				<div class="px-24px pb-20px">
					<p class="text-default mb-12px text-14px">{{ $t('api.directSend.description') }}</p>

					<n-descriptions :column="1" bordered size="small" class="mb-16px">
						<n-descriptions-item :label="$t('api.directSend.fields.endpoint')">
							<code class="text-12px">POST /api/batch_mail/api/send_direct</code>
						</n-descriptions-item>
						<n-descriptions-item :label="$t('api.directSend.fields.auth')">
							<code class="text-12px">X-API-Key: &lt;your-key&gt;</code>
						</n-descriptions-item>
					</n-descriptions>

					<n-table :bordered="true" :single-line="false" size="small" class="mb-16px">
						<thead>
							<tr>
								<th>{{ $t('api.directSend.table.field') }}</th>
								<th>{{ $t('api.directSend.table.required') }}</th>
								<th>{{ $t('api.directSend.table.description') }}</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td><code>to</code></td>
								<td>✅</td>
								<td>{{ $t('api.directSend.table.rows.to') }}</td>
							</tr>
							<tr>
								<td><code>subject</code></td>
								<td>✅</td>
								<td>{{ $t('api.directSend.table.rows.subject') }}</td>
							</tr>
							<tr>
								<td><code>html</code></td>
								<td>✅</td>
								<td>{{ $t('api.directSend.table.rows.html') }}</td>
							</tr>
							<tr>
								<td><code>from</code></td>
								<td>❌</td>
								<td>{{ $t('api.directSend.table.rows.from') }}</td>
							</tr>
							<tr>
								<td><code>from_name</code></td>
								<td>❌</td>
								<td>{{ $t('api.directSend.table.rows.from_name') }}</td>
							</tr>
						</tbody>
					</n-table>

					<p class="text-default mb-8px text-13px font-medium">{{ $t('api.directSend.exampleTitle') }}</p>
					<bt-code :code="exampleCode" language="bash"></bt-code>
				</div>
			</n-collapse-item>
		</n-collapse>
	</n-card>
</template>

<script lang="ts" setup>
const exampleCode = `curl -X POST 'https://your-domain/api/batch_mail/api/send_direct' \\
  -H 'X-API-Key: YOUR_API_KEY' \\
  -H 'Content-Type: application/json' \\
  -d '{
    "to": "user@example.com",
    "subject": "Your verification code",
    "html": "<p>Your code is <strong>123456</strong>.</p>"
  }'`
</script>
