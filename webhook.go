package eleven

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"time"
)

const (
	WebhookSignatureHeader    = "ElevenLabs-Signature"
	webhookTimestampTolerance = 30 * time.Minute
	webhookSignatureScheme    = "v0"
)

type WebhookAuthType string

const (
	WebhookAuthHMAC   WebhookAuthType = "hmac"
	WebhookAuthOAuth2 WebhookAuthType = "oauth2"
	WebhookAuthMTLS   WebhookAuthType = "mtls"
)

func (t WebhookAuthType) String() string { return string(t) }

type WebhookUsageType string

const (
	WebhookUsageConvAIAgentSettings        WebhookUsageType = "ConvAI Agent Settings"
	WebhookUsageConvAISettings             WebhookUsageType = "ConvAI Settings"
	WebhookUsageVoiceLibraryRemovalNotices WebhookUsageType = "Voice Library Removal Notices"
	WebhookUsageSpeechToText               WebhookUsageType = "Speech to Text"
	WebhookUsageAgentQAEvaluations         WebhookUsageType = "Agent QA Evaluations"
	WebhookUsageConvAIAlerting             WebhookUsageType = "ConvAI Alerting"
	WebhookUsageFlows                      WebhookUsageType = "Flows"
	WebhookUsageDubbing                    WebhookUsageType = "Dubbing"
)

func (t WebhookUsageType) String() string { return string(t) }

type WebhookSubscriptionEvent string

const (
	WebhookSubscriptionVoiceLibraryRemovalNotice WebhookSubscriptionEvent = "voice_library_removal_notice"
	WebhookSubscriptionSpeechToText              WebhookSubscriptionEvent = "speech_to_text"
	WebhookSubscriptionAgentQA                   WebhookSubscriptionEvent = "agent_qa"
	WebhookSubscriptionFlows                     WebhookSubscriptionEvent = "flows"
)

func (e WebhookSubscriptionEvent) String() string { return string(e) }

type WebhookEventType string

const (
	WebhookEventPostCallTranscription       WebhookEventType = "post_call_transcription"
	WebhookEventPostCallAudio               WebhookEventType = "post_call_audio"
	WebhookEventCallInitiationFailure       WebhookEventType = "call_initiation_failure"
	WebhookEventVoiceRemovalNotice          WebhookEventType = "voice_removal_notice"
	WebhookEventVoiceRemovalNoticeWithdrawn WebhookEventType = "voice_removal_notice_withdrawn"
	WebhookEventVoiceRemoved                WebhookEventType = "voice_removed"
	WebhookEventFlowsGeneration             WebhookEventType = "flows_generation"
	WebhookEventFlowsTemplateRun            WebhookEventType = "flows_template_run"
	WebhookEventSpeechToTextTranscription   WebhookEventType = "speech_to_text_transcription"
)

func (t WebhookEventType) String() string { return string(t) }

type WebhookUsage struct {
	UsageType WebhookUsageType `json:"usage_type"`
}

type Webhook struct {
	Name                       string                     `json:"name"`
	WebhookID                  string                     `json:"webhook_id"`
	WebhookURL                 string                     `json:"webhook_url"`
	IsDisabled                 bool                       `json:"is_disabled"`
	IsAutoDisabled             bool                       `json:"is_auto_disabled"`
	CreatedAtUnix              int64                      `json:"created_at_unix"`
	AuthType                   WebhookAuthType            `json:"auth_type"`
	Usage                      []WebhookUsage             `json:"usage"`
	Events                     []WebhookSubscriptionEvent `json:"events"`
	MostRecentFailureErrorCode *int                       `json:"most_recent_failure_error_code"`
	MostRecentFailureTimestamp *int64                     `json:"most_recent_failure_timestamp"`
}

type Webhooks []Webhook

type WebhookList struct {
	Webhooks Webhooks `json:"webhooks"`
}

type CreateWebhookResult struct {
	WebhookID     string `json:"webhook_id"`
	WebhookSecret string `json:"webhook_secret"`
}

type WebhookStatusResult struct {
	Status string `json:"status"`
}

type WebhookFilterFunc func(*Webhook) bool

func (w Webhooks) All() iter.Seq[Webhook] {
	return func(yield func(Webhook) bool) {
		for _, webhook := range w {
			if !yield(webhook) {
				return
			}
		}
	}
}

func (w Webhooks) AllWith(filters ...WebhookFilterFunc) iter.Seq[Webhook] {
	return func(yield func(Webhook) bool) {
		for _, webhook := range w {
			if !matchAllWebhook(&webhook, filters) {
				continue
			}
			if !yield(webhook) {
				return
			}
		}
	}
}

func matchAllWebhook(webhook *Webhook, filters []WebhookFilterFunc) bool {
	for _, filter := range filters {
		if !filter(webhook) {
			return false
		}
	}
	return true
}

type WebhookEvent struct {
	Type           WebhookEventType `json:"type"`
	EventTimestamp int64            `json:"event_timestamp"`
	Data           json.RawMessage  `json:"data"`
}

func (e WebhookEvent) UnmarshalData(v any) error {
	if len(e.Data) == 0 {
		return fmt.Errorf("webhook event has no data")
	}
	return json.Unmarshal(e.Data, v)
}

type ConversationTurn struct {
	Role                    string          `json:"role"`
	Message                 string          `json:"message"`
	ToolCalls               json.RawMessage `json:"tool_calls"`
	ToolResults             json.RawMessage `json:"tool_results"`
	Feedback                json.RawMessage `json:"feedback"`
	TimeInCallSecs          int             `json:"time_in_call_secs"`
	ConversationTurnMetrics json.RawMessage `json:"conversation_turn_metrics"`
}

type ConversationMetadata struct {
	StartTimeUnixSecs   int64           `json:"start_time_unix_secs"`
	CallDurationSecs    int             `json:"call_duration_secs"`
	Cost                int             `json:"cost"`
	DeletionSettings    json.RawMessage `json:"deletion_settings"`
	Feedback            json.RawMessage `json:"feedback"`
	AuthorizationMethod string          `json:"authorization_method"`
	Charging            json.RawMessage `json:"charging"`
	TerminationReason   string          `json:"termination_reason"`
}

type ConversationAnalysis struct {
	EvaluationCriteriaResults json.RawMessage `json:"evaluation_criteria_results"`
	DataCollectionResults     json.RawMessage `json:"data_collection_results"`
	CallSuccessful            string          `json:"call_successful"`
	TranscriptSummary         string          `json:"transcript_summary"`
}

type PostCallTranscriptionData struct {
	AgentID                          string               `json:"agent_id"`
	AgentName                        string               `json:"agent_name"`
	ConversationID                   string               `json:"conversation_id"`
	Status                           string               `json:"status"`
	UserID                           string               `json:"user_id"`
	BranchID                         string               `json:"branch_id"`
	VersionID                        string               `json:"version_id"`
	Environment                      string               `json:"environment"`
	Transcript                       []ConversationTurn   `json:"transcript"`
	Metadata                         ConversationMetadata `json:"metadata"`
	Analysis                         ConversationAnalysis `json:"analysis"`
	ConversationInitiationClientData json.RawMessage      `json:"conversation_initiation_client_data"`
	HasAudio                         bool                 `json:"has_audio"`
	HasUserAudio                     bool                 `json:"has_user_audio"`
	HasResponseAudio                 bool                 `json:"has_response_audio"`
}

type PostCallAudioData struct {
	AgentID        string `json:"agent_id"`
	ConversationID string `json:"conversation_id"`
	FullAudio      string `json:"full_audio"`
}

type CallInitiationFailureData struct {
	AgentID        string          `json:"agent_id"`
	ConversationID string          `json:"conversation_id"`
	FailureReason  string          `json:"failure_reason"`
	Metadata       json.RawMessage `json:"metadata"`
}

type FlowsGenerationData struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	ContentURL      string `json:"content_url"`
	ContentMIMEType string `json:"content_mime_type"`
}

type SpeechToTextTranscriptionData struct {
	RequestID       string          `json:"request_id"`
	WebhookMetadata json.RawMessage `json:"webhook_metadata"`
	Transcription   json.RawMessage `json:"transcription"`
}

// ConstructWebhookEvent verifies the ElevenLabs-Signature header, rejects timestamps older than 30 minutes, and unmarshals the raw JSON payload. payload must be the exact request body bytes.
func ConstructWebhookEvent(payload []byte, sigHeader, secret string) (WebhookEvent, error) {
	return constructWebhookEvent(payload, sigHeader, secret, time.Now())
}

func constructWebhookEvent(payload []byte, sigHeader, secret string, now time.Time) (WebhookEvent, error) {
	if sigHeader == "" {
		return WebhookEvent{}, fmt.Errorf("missing signature header")
	}
	if secret == "" {
		return WebhookEvent{}, fmt.Errorf("webhook secret not configured")
	}

	var timestamp, signature string
	for _, part := range strings.Split(sigHeader, ",") {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "t="):
			timestamp = strings.TrimPrefix(part, "t=")
		case strings.HasPrefix(part, webhookSignatureScheme+"="):
			signature = part
		}
	}
	if timestamp == "" || signature == "" {
		return WebhookEvent{}, fmt.Errorf("no signature hash found with expected scheme %s", webhookSignatureScheme)
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return WebhookEvent{}, fmt.Errorf("invalid webhook timestamp")
	}
	if now.Unix()-ts > int64(webhookTimestampTolerance/time.Second) {
		return WebhookEvent{}, fmt.Errorf("timestamp outside the tolerance zone")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	expected := webhookSignatureScheme + "=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return WebhookEvent{}, fmt.Errorf("signature hash does not match the expected signature hash for payload")
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return WebhookEvent{}, fmt.Errorf("failed to marshal webhook payload: %v", err)
	}
	return event, nil
}
