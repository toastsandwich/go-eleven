package eleven

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func TestListWebhooks(t *testing.T) {
	const body = `{"webhooks":[{"name":"My Webhook","webhook_id":"wh_1","webhook_url":"https://example.com/hook","is_disabled":false,"is_auto_disabled":false,"created_at_unix":123,"auth_type":"hmac","usage":[{"usage_type":"Flows"}],"events":["flows"]}]}`
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/workspace/webhooks" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("include_usages") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		io.WriteString(w, body)
	})

	result, err := client.ListWebhooks(WithIncludeUsages(true))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Webhooks) != 1 || result.Webhooks[0].WebhookID != "wh_1" {
		t.Fatalf("webhooks = %+v", result.Webhooks)
	}
	if result.Webhooks[0].AuthType != WebhookAuthHMAC {
		t.Fatalf("auth = %s", result.Webhooks[0].AuthType)
	}
}

func TestCreateWebhook(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/workspace/webhooks" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload createWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Settings.AuthType != "hmac" || payload.Settings.Name != "My Webhook" {
			t.Fatalf("settings = %+v", payload.Settings)
		}
		io.WriteString(w, `{"webhook_id":"wh_1","webhook_secret":"whsec_abc"}`)
	})

	created, err := client.CreateWebhook("My Webhook", "https://example.com/webhook")
	if err != nil {
		t.Fatal(err)
	}
	if created.WebhookID != "wh_1" || created.WebhookSecret != "whsec_abc" {
		t.Fatalf("created = %+v", created)
	}
}

func TestUpdateWebhook(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/workspace/webhooks/wh_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload updateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Name != "My Webhook" || payload.IsDisabled || payload.RetryEnabled == nil || !*payload.RetryEnabled {
			t.Fatalf("payload = %+v", payload)
		}
		io.WriteString(w, `{"status":"ok"}`)
	})

	result, err := client.UpdateWebhook("wh_1", "My Webhook", false, WithWebhookRetryEnabled(true))
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestDeleteWebhook(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/workspace/webhooks/wh_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		io.WriteString(w, `{"status":"ok"}`)
	})

	result, err := client.DeleteWebhook("wh_1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatalf("status = %q", result.Status)
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
	result, err := liveClient(t).ListWebhooks(WithIncludeUsages(true))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("webhooks=%d", len(result.Webhooks))
	for webhook := range result.Webhooks.All() {
		t.Logf("%s id=%s disabled=%v auto=%v auth=%s events=%v usages=%d", webhook.Name, webhook.WebhookID, webhook.IsDisabled, webhook.IsAutoDisabled, webhook.AuthType, webhook.Events, len(webhook.Usage))
	}
}
