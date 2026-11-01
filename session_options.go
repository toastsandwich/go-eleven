package eleven

import (
	"encoding/json"
	"net/url"
)

const defaultSessionEventBuffer = 32

type startConversationOption struct {
	URL              string
	SignedURL        bool
	BranchID         string
	VersionID        string
	UserID           string
	ConfigOverride   json.RawMessage
	DynamicVariables json.RawMessage
	ExtraBody        json.RawMessage
	EventBuffer      int
}

type StartConversationOption func(*startConversationOption)

// Dial this signed WebSocket URL instead of building one from agent_id. Skips xi-api-key.
func WithSessionURL(rawURL string) StartConversationOption {
	return func(o *startConversationOption) {
		o.URL = rawURL
	}
}

// Fetch a signed URL with GetSignedURL, then dial it. Skips xi-api-key. WithSessionURL wins if both are set.
func WithSessionSignedURL() StartConversationOption {
	return func(o *startConversationOption) {
		o.SignedURL = true
	}
}

// The ID of the branch to use.
func WithSessionBranchID(id string) StartConversationOption {
	return func(o *startConversationOption) {
		o.BranchID = id
	}
}

// The ID of the version to use.
func WithSessionVersionID(id string) StartConversationOption {
	return func(o *startConversationOption) {
		o.VersionID = id
	}
}

// Custom user identifier sent in conversation_initiation_client_data.
func WithSessionUserID(id string) StartConversationOption {
	return func(o *startConversationOption) {
		o.UserID = id
	}
}

// Override conversation config for this session as pass-through JSON.
func WithSessionConfigOverride(config json.RawMessage) StartConversationOption {
	return func(o *startConversationOption) {
		o.ConfigOverride = config
	}
}

// Dynamic variables for this session as pass-through JSON.
func WithSessionDynamicVariables(vars json.RawMessage) StartConversationOption {
	return func(o *startConversationOption) {
		o.DynamicVariables = vars
	}
}

// Extra body for a custom LLM as pass-through JSON (custom_llm_extra_body).
func WithSessionExtraBody(body json.RawMessage) StartConversationOption {
	return func(o *startConversationOption) {
		o.ExtraBody = body
	}
}

// Size of the Events channel buffer. Defaults to 32. Values <= 0 use the default.
func WithSessionEventBuffer(n int) StartConversationOption {
	return func(o *startConversationOption) {
		o.EventBuffer = n
	}
}

func (o *startConversationOption) query(agentID string) url.Values {
	q := url.Values{}
	if agentID != "" {
		q.Set("agent_id", agentID)
	}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.VersionID != "" {
		q.Set("version_id", o.VersionID)
	}
	q.Set("source", conversationWSSource)
	return q
}

type conversationInitiationClientData struct {
	Type                       string          `json:"type"`
	UserID                     string          `json:"user_id,omitempty"`
	ConversationConfigOverride json.RawMessage `json:"conversation_config_override,omitempty"`
	DynamicVariables           json.RawMessage `json:"dynamic_variables,omitempty"`
	CustomLLMExtraBody         json.RawMessage `json:"custom_llm_extra_body,omitempty"`
}

func (o *startConversationOption) initiation() conversationInitiationClientData {
	return conversationInitiationClientData{
		Type:                       "conversation_initiation_client_data",
		UserID:                     o.UserID,
		ConversationConfigOverride: o.ConfigOverride,
		DynamicVariables:           o.DynamicVariables,
		CustomLLMExtraBody:         o.ExtraBody,
	}
}

func (o *startConversationOption) eventBuffer() int {
	if o.EventBuffer <= 0 {
		return defaultSessionEventBuffer
	}
	return o.EventBuffer
}
