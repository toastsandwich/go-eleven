package eleven

import (
	"encoding/json"
	"iter"
)

type ToolType string

const (
	ToolTypeWebhook               ToolType = "webhook"
	ToolTypeClient                ToolType = "client"
	ToolTypeSystem                ToolType = "system"
	ToolTypeMCP                   ToolType = "mcp"
	ToolTypeAPIIntegrationWebhook ToolType = "api_integration_webhook"
)

func (t ToolType) String() string { return string(t) }

type ToolSort string

const (
	ToolSortName      ToolSort = "name"
	ToolSortCreatedAt ToolSort = "created_at"
)

func (s ToolSort) String() string { return string(s) }

type ToolExecutionMode string

const (
	ToolExecutionImmediate      ToolExecutionMode = "immediate"
	ToolExecutionPostToolSpeech ToolExecutionMode = "post_tool_speech"
	ToolExecutionAsync          ToolExecutionMode = "async"
)

func (m ToolExecutionMode) String() string { return string(m) }

type ToolInterruptionMode string

const (
	ToolInterruptionAllow                    ToolInterruptionMode = "allow"
	ToolInterruptionDisableDuringTool        ToolInterruptionMode = "disable_during_tool"
	ToolInterruptionDisableDuringToolAndTurn ToolInterruptionMode = "disable_during_tool_and_turn"
)

func (m ToolInterruptionMode) String() string { return string(m) }

type ToolPreSpeech string

const (
	ToolPreSpeechAuto  ToolPreSpeech = "auto"
	ToolPreSpeechForce ToolPreSpeech = "force"
	ToolPreSpeechOff   ToolPreSpeech = "off"
)

func (s ToolPreSpeech) String() string { return string(s) }

type ToolCallSound string

const (
	ToolCallSoundTyping    ToolCallSound = "typing"
	ToolCallSoundElevator1 ToolCallSound = "elevator1"
	ToolCallSoundElevator2 ToolCallSound = "elevator2"
	ToolCallSoundElevator3 ToolCallSound = "elevator3"
	ToolCallSoundElevator4 ToolCallSound = "elevator4"
)

func (s ToolCallSound) String() string { return string(s) }

type ToolCallSoundBehavior string

const (
	ToolCallSoundAuto   ToolCallSoundBehavior = "auto"
	ToolCallSoundAlways ToolCallSoundBehavior = "always"
)

func (b ToolCallSoundBehavior) String() string { return string(b) }

type ToolErrorHandlingMode string

const (
	ToolErrorHandlingAuto        ToolErrorHandlingMode = "auto"
	ToolErrorHandlingSummarized  ToolErrorHandlingMode = "summarized"
	ToolErrorHandlingPassthrough ToolErrorHandlingMode = "passthrough"
	ToolErrorHandlingHide        ToolErrorHandlingMode = "hide"
)

func (m ToolErrorHandlingMode) String() string { return string(m) }

type DependentAgentType string

const (
	DependentAgentAvailable DependentAgentType = "available"
	DependentAgentUnknown   DependentAgentType = "unknown"
)

func (t DependentAgentType) String() string { return string(t) }

type WebhookAPISchema struct {
	URL               string          `json:"url"`
	Method            string          `json:"method,omitempty"`
	PathParamsSchema  json.RawMessage `json:"path_params_schema,omitempty"`
	QueryParamsSchema json.RawMessage `json:"query_params_schema,omitempty"`
	RequestBodySchema json.RawMessage `json:"request_body_schema,omitempty"`
	RequestHeaders    json.RawMessage `json:"request_headers,omitempty"`
	AuthConnection    json.RawMessage `json:"auth_connection,omitempty"`
}

type ToolConfig struct {
	Type        ToolType `json:"type"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`

	APISchema                     *WebhookAPISchema `json:"api_schema,omitempty"`
	FollowRedirects               *bool             `json:"follow_redirects,omitempty"`
	FollowRedirectsAllowedDomains []string          `json:"follow_redirects_allowed_domains,omitempty"`

	ExpectsResponse *bool           `json:"expects_response,omitempty"`
	Parameters      json.RawMessage `json:"parameters,omitempty"`

	Params json.RawMessage `json:"params,omitempty"`
	Value  json.RawMessage `json:"value,omitempty"`

	ResponseTimeoutSecs   *int                  `json:"response_timeout_secs,omitempty"`
	ExecutionMode         ToolExecutionMode     `json:"execution_mode,omitempty"`
	InterruptionMode      ToolInterruptionMode  `json:"interruption_mode,omitempty"`
	PreToolSpeech         ToolPreSpeech         `json:"pre_tool_speech,omitempty"`
	ToolCallSound         ToolCallSound         `json:"tool_call_sound,omitempty"`
	ToolCallSoundBehavior ToolCallSoundBehavior `json:"tool_call_sound_behavior,omitempty"`
	ToolErrorHandlingMode ToolErrorHandlingMode `json:"tool_error_handling_mode,omitempty"`
	Assignments           json.RawMessage       `json:"assignments,omitempty"`
	DynamicVariables      json.RawMessage       `json:"dynamic_variables,omitempty"`
}

type ToolUsageStats struct {
	AvgLatencySecs float64 `json:"avg_latency_secs"`
	TotalCalls     int     `json:"total_calls"`
}

type Tool struct {
	ID            string             `json:"id"`
	ToolConfig    ToolConfig         `json:"tool_config"`
	AccessInfo    ResourceAccessInfo `json:"access_info"`
	UsageStats    ToolUsageStats     `json:"usage_stats"`
	ResponseMocks json.RawMessage    `json:"response_mocks,omitempty"`
}

type Tools []Tool

type ToolFilterFunc func(*Tool) bool

func (t Tools) All() iter.Seq[Tool] {
	return func(yield func(Tool) bool) {
		for _, tool := range t {
			if !yield(tool) {
				return
			}
		}
	}
}

func (t Tools) AllWith(filters ...ToolFilterFunc) iter.Seq[Tool] {
	return func(yield func(Tool) bool) {
		for _, tool := range t {
			if !matchAllTool(&tool, filters) {
				continue
			}
			if !yield(tool) {
				return
			}
		}
	}
}

func matchAllTool(tool *Tool, filters []ToolFilterFunc) bool {
	for _, filter := range filters {
		if !filter(tool) {
			return false
		}
	}
	return true
}

type ToolList struct {
	Tools      Tools  `json:"tools"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type DependentAgent struct {
	Type                  DependentAgentType `json:"type"`
	ID                    string             `json:"id"`
	Name                  string             `json:"name,omitempty"`
	AccessLevel           string             `json:"access_level,omitempty"`
	CreatedAtUnixSecs     int64              `json:"created_at_unix_secs,omitempty"`
	ReferencedResourceIDs []string           `json:"referenced_resource_ids,omitempty"`
}

type DependentBranch struct {
	AgentID    string `json:"agent_id"`
	AgentName  string `json:"agent_name"`
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	IsMain     bool   `json:"is_main"`
}

type DependentAgentList struct {
	Agents     []DependentAgent  `json:"agents"`
	Branches   []DependentBranch `json:"branches,omitempty"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}
