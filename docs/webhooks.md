# Webhooks

Create, list, update, and delete workspace webhooks with `/v1/workspace/webhooks`. Verify incoming deliveries with `ConstructWebhookEvent`.

[Back to docs](README.md)

## List webhooks

```go
result, err := client.ListWebhooks(eleven.WithIncludeUsages(true))
if err != nil {
	log.Fatal(err)
}

for webhook := range result.Webhooks.All() {
	fmt.Println(webhook.Name, webhook.WebhookID, webhook.WebhookURL)
}
```

`Webhooks` also has `AllWith` for extra local filters after the API response. `usage` and `events` are populated when usages are requested (workspace admins only).

## Create webhook

Creates an HMAC-authenticated webhook. Store `WebhookSecret` and use it to verify deliveries.

```go
created, err := client.CreateWebhook(
	"My Webhook",
	"https://api.example.com/webhook",
	eleven.WithWebhookRequestHeaders(map[string]string{
		"Authorization": "Bearer your-token",
	}),
)
if err != nil {
	log.Fatal(err)
}
fmt.Println(created.WebhookID, created.WebhookSecret)
```

Custom headers are write-only. `ListWebhooks` does not return them.

## Update webhook

`name` and `isDisabled` are required. Pass the current name if you are only changing retries, headers, or events.

```go
_, err := client.UpdateWebhook(
	webhookID,
	"My Webhook",
	false,
	eleven.WithWebhookRetryEnabled(true),
	eleven.WithWebhookEvents(eleven.WebhookSubscriptionFlows, eleven.WebhookSubscriptionSpeechToText),
)
```

| Option | Purpose |
| --- | --- |
| `WithWebhookRetryEnabled` | Retry 5xx, 429, and timeouts (up to 5 attempts) |
| `WithUpdateWebhookRequestHeaders` | Replace delivery headers |
| `WithWebhookEvents` | Replace workspace-level subscriptions |

Workspace-level events: `voice_library_removal_notice`, `speech_to_text`, `agent_qa`, `flows`. Omit `WithWebhookEvents` to leave subscriptions unchanged. Passing no events unsubscribes from all of them.

## Delete webhook

```go
_, err := client.DeleteWebhook(webhookID)
```

The webhook must not be in use by any product.

## Verify incoming events

Use the raw request body. Parsing JSON first and re-encoding it will fail HMAC verification.

```go
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	event, err := eleven.ConstructWebhookEvent(
		body,
		r.Header.Get(eleven.WebhookSignatureHeader),
		os.Getenv("WEBHOOK_SECRET"),
	)
	if err != nil {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	switch event.Type {
	case eleven.WebhookEventPostCallTranscription:
		var data eleven.PostCallTranscriptionData
		if err := event.UnmarshalData(&data); err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}
		fmt.Println(data.ConversationID, data.Status)
	case eleven.WebhookEventPostCallAudio:
		var data eleven.PostCallAudioData
		if err := event.UnmarshalData(&data); err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}
		fmt.Println(data.ConversationID)
	case eleven.WebhookEventFlowsGeneration:
		var data eleven.FlowsGenerationData
		if err := event.UnmarshalData(&data); err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}
		fmt.Println(data.ID, data.Status, data.ContentURL)
	}

	w.WriteHeader(http.StatusOK)
})
```

`ConstructWebhookEvent` checks the `ElevenLabs-Signature` header (`t=<unix>,v0=<hmac-sha256>`), rejects timestamps older than 30 minutes, and unmarshals the payload. Return HTTP 200 quickly after verification. Handlers should be idempotent; retries send the same payload.

## Event types

| `event.Type` | Data |
| --- | --- |
| `post_call_transcription` | `PostCallTranscriptionData` |
| `post_call_audio` | `PostCallAudioData` |
| `call_initiation_failure` | `CallInitiationFailureData` |
| `flows_generation` | `FlowsGenerationData` |
| `speech_to_text_transcription` | `SpeechToTextTranscriptionData` |
| `voice_removal_notice` | unmarshal into your own struct |
| `voice_removal_notice_withdrawn` | unmarshal into your own struct |
| `voice_removed` | unmarshal into your own struct |
| `flows_template_run` | unmarshal into your own struct |

## Context

```go
result, err := client.ListWebhooksWithContext(ctx, eleven.WithIncludeUsages(true))
```
