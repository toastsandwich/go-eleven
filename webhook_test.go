package eleven

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func signWebhook(payload, secret string, ts int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.", ts)))
	mac.Write([]byte(payload))
	return fmt.Sprintf("t=%d,v0=%s", ts, hex.EncodeToString(mac.Sum(nil)))
}

func TestConstructWebhookEvent(t *testing.T) {
	const secret = "whsec_test"
	payload := []byte(`{"type":"post_call_transcription","event_timestamp":1739537297,"data":{"agent_id":"xyz","conversation_id":"abc","status":"done"}}`)
	now := time.Unix(1739537297, 0)

	t.Run("valid", func(t *testing.T) {
		header := signWebhook(string(payload), secret, now.Unix())
		event, err := constructWebhookEvent(payload, header, secret, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if event.Type != WebhookEventPostCallTranscription {
			t.Fatalf("type = %q", event.Type)
		}
		if event.EventTimestamp != 1739537297 {
			t.Fatalf("event_timestamp = %d", event.EventTimestamp)
		}

		var data PostCallTranscriptionData
		if err := event.UnmarshalData(&data); err != nil {
			t.Fatalf("unmarshal data: %v", err)
		}
		if data.AgentID != "xyz" || data.ConversationID != "abc" || data.Status != "done" {
			t.Fatalf("data = %+v", data)
		}
	})

	t.Run("missing signature header", func(t *testing.T) {
		_, err := constructWebhookEvent(payload, "", secret, now)
		if err == nil || err.Error() != "missing signature header" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("missing secret", func(t *testing.T) {
		header := signWebhook(string(payload), secret, now.Unix())
		_, err := constructWebhookEvent(payload, header, "", now)
		if err == nil || err.Error() != "webhook secret not configured" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("missing scheme", func(t *testing.T) {
		_, err := constructWebhookEvent(payload, fmt.Sprintf("t=%d", now.Unix()), secret, now)
		if err == nil || err.Error() != "no signature hash found with expected scheme v0" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("expired timestamp", func(t *testing.T) {
		ts := now.Add(-31 * time.Minute).Unix()
		header := signWebhook(string(payload), secret, ts)
		_, err := constructWebhookEvent(payload, header, secret, now)
		if err == nil || err.Error() != "timestamp outside the tolerance zone" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("timestamp at tolerance", func(t *testing.T) {
		ts := now.Add(-webhookTimestampTolerance).Unix()
		header := signWebhook(string(payload), secret, ts)
		if _, err := constructWebhookEvent(payload, header, secret, now); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		header := signWebhook(string(payload), "other", now.Unix())
		_, err := constructWebhookEvent(payload, header, secret, now)
		if err == nil || err.Error() != "signature hash does not match the expected signature hash for payload" {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("tampered payload", func(t *testing.T) {
		header := signWebhook(string(payload), secret, now.Unix())
		tampered := []byte(`{"type":"post_call_audio","event_timestamp":1739537297,"data":{}}`)
		_, err := constructWebhookEvent(tampered, header, secret, now)
		if err == nil || err.Error() != "signature hash does not match the expected signature hash for payload" {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestWebhooksAllWith(t *testing.T) {
	webhooks := Webhooks{
		{Name: "a", IsDisabled: true},
		{Name: "b", IsDisabled: false},
		{Name: "c", IsDisabled: false},
	}

	var names []string
	for webhook := range webhooks.AllWith(func(w *Webhook) bool {
		return !w.IsDisabled
	}) {
		names = append(names, webhook.Name)
	}
	if len(names) != 2 || names[0] != "b" || names[1] != "c" {
		t.Fatalf("names = %v", names)
	}
}

func TestCreateWebhookValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.CreateWebhook("", "https://example.com/webhook"); err == nil || err.Error() != "name is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.CreateWebhook("hook", ""); err == nil || err.Error() != "webhookURL is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestUpdateWebhookValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.UpdateWebhook("", "hook", false); err == nil || err.Error() != "webhookID is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.UpdateWebhook("id", "", false); err == nil || err.Error() != "name is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestDeleteWebhookValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.DeleteWebhook(""); err == nil || err.Error() != "webhookID is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestCreateWebhookRequestBody(t *testing.T) {
	var opt createWebhookOption
	WithWebhookRequestHeaders(map[string]string{"Authorization": "Bearer tok"})(&opt)
	got, err := json.Marshal(opt.body("My Webhook", "https://example.com/webhook"))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"settings":{"auth_type":"hmac","name":"My Webhook","webhook_url":"https://example.com/webhook","request_headers":{"Authorization":"Bearer tok"}}}`
	if string(got) != want {
		t.Fatalf("got %s", got)
	}
}

func TestUpdateWebhookRequestBody(t *testing.T) {
	var opt updateWebhookOption
	WithWebhookRetryEnabled(true)(&opt)
	WithWebhookEvents()(&opt)
	got, err := json.Marshal(opt.body("My Webhook", false))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"is_disabled":false,"name":"My Webhook","retry_enabled":true,"events":[]}`
	if string(got) != want {
		t.Fatalf("got %s", got)
	}
}

func TestListWebhooksQuery(t *testing.T) {
	var opt listWebhooksOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}
	WithIncludeUsages(true)(&opt)
	if q := opt.query().Encode(); q != "include_usages=true" {
		t.Fatalf("query = %q", q)
	}
}

func TestLiveListWebhooks(t *testing.T) {
	if os.Getenv("ELEVEN_LIVE") == "" {
		t.Skip("set ELEVEN_LIVE=1 to hit the API")
	}
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		t.Fatal("ELEVEN_API_KEY is empty")
	}
	result, err := NewClient(apiKey).ListWebhooks(WithIncludeUsages(true))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("webhooks=%d", len(result.Webhooks))
	for webhook := range result.Webhooks.All() {
		t.Logf("%s id=%s disabled=%v auto=%v auth=%s events=%v usages=%d", webhook.Name, webhook.WebhookID, webhook.IsDisabled, webhook.IsAutoDisabled, webhook.AuthType, webhook.Events, len(webhook.Usage))
	}
}
