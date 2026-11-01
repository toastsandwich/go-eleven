package eleven

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type listConversationsOption struct {
	Cursor          string
	AgentID         string
	PageSize        int
	SortDirection   string
	CallSuccessful  string
	CallStartAfter  *int64
	CallStartBefore *int64
	UserID          string
	BranchID        string
	VersionID       string
	SummaryMode     string
}

type ListConversationsOption func(*listConversationsOption)

// Used for fetching the next page. Cursor is returned in the list response as NextCursor.
func WithConversationCursor(cursor string) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.Cursor = cursor
	}
}

// Restrict results to conversations of this agent.
func WithConversationAgentID(id string) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.AgentID = id
	}
}

// How many conversations to return at maximum. Can not exceed 100, defaults to 30.
func WithConversationPageSize(size int) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.PageSize = size
	}
}

// The direction to sort conversations by call start time. Defaults to descending (newest first).
func WithConversationSortDirection(direction SortDirection) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.SortDirection = direction.String()
	}
}

// The result of the success evaluation: success, failure, unknown, or error.
func WithCallSuccessful(result CallSuccessful) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.CallSuccessful = result.String()
	}
}

// Unix timestamp (in seconds) to filter conversations after this start date.
func WithCallStartAfterUnix(ts int64) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.CallStartAfter = &ts
	}
}

// Unix timestamp (in seconds) to filter conversations up to this start date.
func WithCallStartBeforeUnix(ts int64) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.CallStartBefore = &ts
	}
}

// Filter conversations by the user ID who initiated them.
func WithConversationUserID(id string) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.UserID = id
	}
}

// Filter conversations by branch ID.
func WithConversationBranchID(id string) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.BranchID = id
	}
}

// Filter conversations by version ID.
func WithConversationVersionID(id string) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.VersionID = id
	}
}

// Whether to include transcript summaries in the list response. exclude or include. Defaults to exclude.
func WithSummaryMode(mode ConversationSummaryMode) ListConversationsOption {
	return func(o *listConversationsOption) {
		o.SummaryMode = mode.String()
	}
}

func (o *listConversationsOption) query() url.Values {
	q := url.Values{}
	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	if o.AgentID != "" {
		q.Set("agent_id", o.AgentID)
	}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.SortDirection != "" {
		q.Set("sort_direction", o.SortDirection)
	}
	if o.CallSuccessful != "" {
		q.Set("call_successful", o.CallSuccessful)
	}
	if o.CallStartAfter != nil {
		q.Set("call_start_after_unix", strconv.FormatInt(*o.CallStartAfter, 10))
	}
	if o.CallStartBefore != nil {
		q.Set("call_start_before_unix", strconv.FormatInt(*o.CallStartBefore, 10))
	}
	if o.UserID != "" {
		q.Set("user_id", o.UserID)
	}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.VersionID != "" {
		q.Set("version_id", o.VersionID)
	}
	if o.SummaryMode != "" {
		q.Set("summary_mode", o.SummaryMode)
	}
	return q
}

type getConversationSummaryOption struct {
	MaxMessages *int
}

type GetConversationSummaryOption func(*getConversationSummaryOption)

// Maximum number of chat message turns to include inline. When the conversation has more than this, the messages are omitted and messages_omitted is set. Defaults to 40.
func WithMaxMessages(n int) GetConversationSummaryOption {
	return func(o *getConversationSummaryOption) {
		o.MaxMessages = &n
	}
}

func (o *getConversationSummaryOption) query() url.Values {
	q := url.Values{}
	if o.MaxMessages != nil {
		q.Set("max_messages", strconv.Itoa(*o.MaxMessages))
	}
	return q
}

type signedURLOption struct {
	IncludeConversationID *bool
	BranchID              string
	VersionID             string
	Environment           string
}

type SignedURLOption func(*signedURLOption)

// Whether to include a conversation_id with the response. If included, the conversation_signature cannot be used again.
func WithSignedURLConversationID() SignedURLOption {
	return func(o *signedURLOption) {
		include := true
		o.IncludeConversationID = &include
	}
}

// The ID of the branch to use.
func WithSignedURLBranchID(id string) SignedURLOption {
	return func(o *signedURLOption) {
		o.BranchID = id
	}
}

// The ID of the version to use.
func WithSignedURLVersionID(id string) SignedURLOption {
	return func(o *signedURLOption) {
		o.VersionID = id
	}
}

// The environment to use for resolving environment variables (e.g. production, staging). Defaults to production.
func WithSignedURLEnvironment(env string) SignedURLOption {
	return func(o *signedURLOption) {
		o.Environment = env
	}
}

func (o *signedURLOption) query(agentID string) url.Values {
	q := url.Values{}
	q.Set("agent_id", agentID)
	if o.IncludeConversationID != nil {
		q.Set("include_conversation_id", strconv.FormatBool(*o.IncludeConversationID))
	}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.VersionID != "" {
		q.Set("version_id", o.VersionID)
	}
	if o.Environment != "" {
		q.Set("environment", o.Environment)
	}
	return q
}

type conversationTokenOption struct {
	ParticipantName string
	BranchID        string
	VersionID       string
	Environment     string
}

type ConversationTokenOption func(*conversationTokenOption)

// Optional custom participant name. If not provided, user ID will be used.
func WithParticipantName(name string) ConversationTokenOption {
	return func(o *conversationTokenOption) {
		o.ParticipantName = name
	}
}

// The ID of the branch to use.
func WithTokenBranchID(id string) ConversationTokenOption {
	return func(o *conversationTokenOption) {
		o.BranchID = id
	}
}

// The ID of the version to use.
func WithTokenVersionID(id string) ConversationTokenOption {
	return func(o *conversationTokenOption) {
		o.VersionID = id
	}
}

// The environment to use for resolving environment variables (e.g. production, staging). Defaults to production.
func WithTokenEnvironment(env string) ConversationTokenOption {
	return func(o *conversationTokenOption) {
		o.Environment = env
	}
}

func (o *conversationTokenOption) query(agentID string) url.Values {
	q := url.Values{}
	q.Set("agent_id", agentID)
	if o.ParticipantName != "" {
		q.Set("participant_name", o.ParticipantName)
	}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.VersionID != "" {
		q.Set("version_id", o.VersionID)
	}
	if o.Environment != "" {
		q.Set("environment", o.Environment)
	}
	return q
}

type conversationFeedbackRequest struct {
	Feedback string `json:"feedback"`
}

type clientSecretOption struct {
	Session      json.RawMessage
	ExpiresAfter *int
}

type ClientSecretOption func(*clientSecretOption)

// OpenAI-compatible Realtime session config as pass-through JSON.
func WithRealtimeSession(session json.RawMessage) ClientSecretOption {
	return func(o *clientSecretOption) {
		o.Session = session
	}
}

// Lifetime of the client secret in seconds, anchored at created_at.
func WithRealtimeExpiresAfter(seconds int) ClientSecretOption {
	return func(o *clientSecretOption) {
		o.ExpiresAfter = &seconds
	}
}

type realtimeExpiresAfter struct {
	Seconds int    `json:"seconds"`
	Anchor  string `json:"anchor,omitempty"`
}

type createClientSecretRequest struct {
	Session      json.RawMessage       `json:"session,omitempty"`
	ExpiresAfter *realtimeExpiresAfter `json:"expires_after,omitempty"`
}

func (o *clientSecretOption) body() createClientSecretRequest {
	req := createClientSecretRequest{Session: o.Session}
	if o.ExpiresAfter != nil {
		req.ExpiresAfter = &realtimeExpiresAfter{Seconds: *o.ExpiresAfter, Anchor: "created_at"}
	}
	return req
}
