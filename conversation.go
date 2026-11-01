package eleven

import (
	"encoding/json"
	"iter"
)

type ConversationStatus string

const (
	ConversationStatusInitiated  ConversationStatus = "initiated"
	ConversationStatusInProgress ConversationStatus = "in-progress"
	ConversationStatusProcessing ConversationStatus = "processing"
	ConversationStatusDone       ConversationStatus = "done"
	ConversationStatusFailed     ConversationStatus = "failed"
)

func (s ConversationStatus) String() string { return string(s) }

type CallSuccessful string

const (
	CallSuccessfulSuccess CallSuccessful = "success"
	CallSuccessfulFailure CallSuccessful = "failure"
	CallSuccessfulUnknown CallSuccessful = "unknown"
	CallSuccessfulError   CallSuccessful = "error"
)

func (s CallSuccessful) String() string { return string(s) }

type ConversationDirection string

const (
	ConversationDirectionInbound  ConversationDirection = "inbound"
	ConversationDirectionOutbound ConversationDirection = "outbound"
)

func (d ConversationDirection) String() string { return string(d) }

type ConversationFeedback string

const (
	ConversationFeedbackLike    ConversationFeedback = "like"
	ConversationFeedbackDislike ConversationFeedback = "dislike"
)

func (f ConversationFeedback) String() string { return string(f) }

type ConversationSummaryMode string

const (
	ConversationSummaryExclude ConversationSummaryMode = "exclude"
	ConversationSummaryInclude ConversationSummaryMode = "include"
)

func (m ConversationSummaryMode) String() string { return string(m) }

type ConversationListItem struct {
	AgentID                      string                `json:"agent_id"`
	ConversationID               string                `json:"conversation_id"`
	StartTimeUnixSecs            int64                 `json:"start_time_unix_secs"`
	CallDurationSecs             int                   `json:"call_duration_secs"`
	MessageCount                 int                   `json:"message_count"`
	Status                       ConversationStatus    `json:"status"`
	CallSuccessful               CallSuccessful        `json:"call_successful"`
	BranchID                     string                `json:"branch_id,omitempty"`
	VersionID                    string                `json:"version_id,omitempty"`
	AgentName                    string                `json:"agent_name,omitempty"`
	TerminationReason            string                `json:"termination_reason,omitempty"`
	TranscriptSummary            string                `json:"transcript_summary,omitempty"`
	CallSummaryTitle             string                `json:"call_summary_title,omitempty"`
	MainLanguage                 string                `json:"main_language,omitempty"`
	ConversationInitiationSource string                `json:"conversation_initiation_source,omitempty"`
	Direction                    ConversationDirection `json:"direction,omitempty"`
	Rating                       *float64              `json:"rating,omitempty"`
	ToolNames                    []string              `json:"tool_names,omitempty"`
	TagIDs                       []string              `json:"tag_ids,omitempty"`
}

type ConversationListItems []ConversationListItem

type ConversationFilterFunc func(*ConversationListItem) bool

func (c ConversationListItems) All() iter.Seq[ConversationListItem] {
	return func(yield func(ConversationListItem) bool) {
		for _, conversation := range c {
			if !yield(conversation) {
				return
			}
		}
	}
}

func (c ConversationListItems) AllWith(filters ...ConversationFilterFunc) iter.Seq[ConversationListItem] {
	return func(yield func(ConversationListItem) bool) {
		for _, conversation := range c {
			if !matchAllConversation(&conversation, filters) {
				continue
			}
			if !yield(conversation) {
				return
			}
		}
	}
}

func matchAllConversation(conversation *ConversationListItem, filters []ConversationFilterFunc) bool {
	for _, filter := range filters {
		if !filter(conversation) {
			return false
		}
	}
	return true
}

type ConversationList struct {
	Conversations ConversationListItems `json:"conversations"`
	HasMore       bool                  `json:"has_more"`
	NextCursor    string                `json:"next_cursor"`
}

type ConversationHistory struct {
	StartTimeUnixSecs   int64           `json:"start_time_unix_secs"`
	CallDurationSecs    int             `json:"call_duration_secs"`
	Cost                *int            `json:"cost,omitempty"`
	TerminationReason   string          `json:"termination_reason,omitempty"`
	AuthorizationMethod string          `json:"authorization_method,omitempty"`
	MainLanguage        string          `json:"main_language,omitempty"`
	TextOnly            bool            `json:"text_only,omitempty"`
	Feedback            json.RawMessage `json:"feedback,omitempty"`
	Charging            json.RawMessage `json:"charging,omitempty"`
	PhoneCall           json.RawMessage `json:"phone_call,omitempty"`
	Error               json.RawMessage `json:"error,omitempty"`
}

type ConversationTranscriptTurn struct {
	Role           string          `json:"role"`
	Message        string          `json:"message,omitempty"`
	TimeInCallSecs int             `json:"time_in_call_secs"`
	Interrupted    bool            `json:"interrupted,omitempty"`
	ToolCalls      json.RawMessage `json:"tool_calls,omitempty"`
	ToolResults    json.RawMessage `json:"tool_results,omitempty"`
	Feedback       json.RawMessage `json:"feedback,omitempty"`
}

type ConversationCallAnalysis struct {
	CallSuccessful            CallSuccessful  `json:"call_successful"`
	TranscriptSummary         string          `json:"transcript_summary"`
	CallSummaryTitle          string          `json:"call_summary_title,omitempty"`
	EvaluationCriteriaResults json.RawMessage `json:"evaluation_criteria_results,omitempty"`
	DataCollectionResults     json.RawMessage `json:"data_collection_results,omitempty"`
}

type Conversation struct {
	AgentID                          string                       `json:"agent_id"`
	ConversationID                   string                       `json:"conversation_id"`
	Status                           ConversationStatus           `json:"status"`
	Metadata                         ConversationHistory          `json:"metadata"`
	Transcript                       []ConversationTranscriptTurn `json:"transcript"`
	HasAudio                         bool                         `json:"has_audio"`
	HasUserAudio                     bool                         `json:"has_user_audio"`
	HasResponseAudio                 bool                         `json:"has_response_audio"`
	HasAuxiliaryAudio                bool                         `json:"has_auxiliary_audio"`
	AgentName                        string                       `json:"agent_name,omitempty"`
	UserID                           string                       `json:"user_id,omitempty"`
	BranchID                         string                       `json:"branch_id,omitempty"`
	VersionID                        string                       `json:"version_id,omitempty"`
	Environment                      string                       `json:"environment,omitempty"`
	Analysis                         *ConversationCallAnalysis    `json:"analysis,omitempty"`
	ConversationInitiationClientData json.RawMessage              `json:"conversation_initiation_client_data,omitempty"`
	TagIDs                           []string                     `json:"tag_ids,omitempty"`
}

type ConversationSummaryMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type ConversationSummary struct {
	ConversationID    string                       `json:"conversation_id"`
	AgentID           string                       `json:"agent_id"`
	Status            ConversationStatus           `json:"status"`
	MessageCount      int                          `json:"message_count"`
	Note              string                       `json:"note"`
	CallSummaryTitle  string                       `json:"call_summary_title,omitempty"`
	TranscriptSummary string                       `json:"transcript_summary,omitempty"`
	CallSuccessful    CallSuccessful               `json:"call_successful,omitempty"`
	Messages          []ConversationSummaryMessage `json:"messages,omitempty"`
	MessagesOmitted   bool                         `json:"messages_omitted"`
}

type SignedURL struct {
	SignedURL string `json:"signed_url"`
}

type ConversationToken struct {
	Token          string `json:"token"`
	ConversationID string `json:"conversation_id"`
}

type RealtimeClientSecret struct {
	Value     string          `json:"value"`
	ExpiresAt int64           `json:"expires_at"`
	Session   json.RawMessage `json:"session,omitempty"`
}
