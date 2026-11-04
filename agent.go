package eleven

import (
	"encoding/json"
	"iter"
)

type AgentSort string

const (
	AgentSortName        AgentSort = "name"
	AgentSortCreatedAt   AgentSort = "created_at"
	AgentSortUpdatedAt   AgentSort = "updated_at"
	AgentSortCallCount7d AgentSort = "call_count_7d"
)

func (s AgentSort) String() string { return string(s) }

type AgentAccessRole string

const (
	AgentAccessRoleAdmin     AgentAccessRole = "admin"
	AgentAccessRoleEditor    AgentAccessRole = "editor"
	AgentAccessRoleCommenter AgentAccessRole = "commenter"
	AgentAccessRoleViewer    AgentAccessRole = "viewer"
)

func (r AgentAccessRole) String() string { return string(r) }

type AgentAccessSource string

const (
	AgentAccessSourceCreator          AgentAccessSource = "creator"
	AgentAccessSourceExplicit         AgentAccessSource = "explicit"
	AgentAccessSourceWorkspaceAdmin   AgentAccessSource = "workspace_admin"
	AgentAccessSourceWorkspaceDefault AgentAccessSource = "workspace_default"
)

func (s AgentAccessSource) String() string { return string(s) }

type ResourceAccessInfo struct {
	IsCreator                    bool              `json:"is_creator"`
	CreatorName                  string            `json:"creator_name"`
	CreatorEmail                 string            `json:"creator_email"`
	Role                         AgentAccessRole   `json:"role"`
	AnonymousAccessLevelOverride AgentAccessRole   `json:"anonymous_access_level_override,omitempty"`
	AccessSource                 AgentAccessSource `json:"access_source,omitempty"`
}

type AgentSummary struct {
	AgentID              string             `json:"agent_id"`
	Name                 string             `json:"name"`
	VoiceID              string             `json:"voice_id"`
	Tags                 []string           `json:"tags"`
	CreatedAtUnixSecs    int64              `json:"created_at_unix_secs"`
	UpdatedAtUnixSecs    int64              `json:"updated_at_unix_secs"`
	AccessInfo           ResourceAccessInfo `json:"access_info"`
	LastCallTimeUnixSecs *int64             `json:"last_call_time_unix_secs"`
	Archived             bool               `json:"archived"`
}

type Agents []AgentSummary

type AgentFilterFunc func(*AgentSummary) bool

func (a Agents) All() iter.Seq[AgentSummary] {
	return func(yield func(AgentSummary) bool) {
		for _, agent := range a {
			if !yield(agent) {
				return
			}
		}
	}
}

func (a Agents) AllWith(filters ...AgentFilterFunc) iter.Seq[AgentSummary] {
	return func(yield func(AgentSummary) bool) {
		for _, agent := range a {
			if !matchAllAgent(&agent, filters) {
				continue
			}
			if !yield(agent) {
				return
			}
		}
	}
}

func matchAllAgent(agent *AgentSummary, filters []AgentFilterFunc) bool {
	for _, filter := range filters {
		if !filter(agent) {
			return false
		}
	}
	return true
}

type AgentList struct {
	Agents     Agents `json:"agents"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type AgentMetadata struct {
	CreatedAtUnixSecs int64 `json:"created_at_unix_secs"`
	UpdatedAtUnixSecs int64 `json:"updated_at_unix_secs"`
}

type ASRConfig struct {
	Quality              string   `json:"quality,omitempty"`
	Provider             string   `json:"provider,omitempty"`
	UserInputAudioFormat string   `json:"user_input_audio_format,omitempty"`
	Keywords             []string `json:"keywords,omitempty"`
}

type SoftTimeoutConfig struct {
	TimeoutSeconds                    *float64 `json:"timeout_seconds,omitempty"`
	UseLLMGeneratedMessage            *bool    `json:"use_llm_generated_message,omitempty"`
	LLMGeneratedMessagePromptOverride string   `json:"llm_generated_message_prompt_override,omitempty"`
	MaxSoftTimeoutsPerGeneration      *int     `json:"max_soft_timeouts_per_generation,omitempty"`
}

type TurnConfig struct {
	TurnTimeout           *float64           `json:"turn_timeout,omitempty"`
	SilenceEndCallTimeout *float64           `json:"silence_end_call_timeout,omitempty"`
	TurnEagerness         string             `json:"turn_eagerness,omitempty"`
	SpeculativeTurn       *bool              `json:"speculative_turn,omitempty"`
	SoftTimeoutConfig     *SoftTimeoutConfig `json:"soft_timeout_config,omitempty"`
}

type AgentTTSConfig struct {
	ModelID                string   `json:"model_id,omitempty"`
	VoiceID                string   `json:"voice_id,omitempty"`
	AgentOutputAudioFormat string   `json:"agent_output_audio_format,omitempty"`
	Stability              *float64 `json:"stability,omitempty"`
	Speed                  *float64 `json:"speed,omitempty"`
	SimilarityBoost        *float64 `json:"similarity_boost,omitempty"`
	ExpressiveMode         *bool    `json:"expressive_mode,omitempty"`
}

type AgentConversation struct {
	TextOnly           *bool `json:"text_only,omitempty"`
	MaxDurationSeconds *int  `json:"max_duration_seconds,omitempty"`
}

type AgentKnowledgeBaseLocator struct {
	Type      KnowledgeType  `json:"type"`
	Name      string         `json:"name"`
	ID        string         `json:"id"`
	UsageMode KnowledgeUsage `json:"usage_mode,omitempty"`
}

type AgentPrompt struct {
	Prompt        string                      `json:"prompt,omitempty"`
	LLM           string                      `json:"llm,omitempty"`
	Temperature   *float64                    `json:"temperature,omitempty"`
	MaxTokens     *int                        `json:"max_tokens,omitempty"`
	ToolIDs       []string                    `json:"tool_ids,omitempty"`
	MCPServerIDs  []string                    `json:"mcp_server_ids,omitempty"`
	KnowledgeBase []AgentKnowledgeBaseLocator `json:"knowledge_base,omitempty"`
	RAG           json.RawMessage             `json:"rag,omitempty"`
}

type AgentConfig struct {
	FirstMessage                string            `json:"first_message,omitempty"`
	Language                    string            `json:"language,omitempty"`
	Prompt                      *AgentPrompt      `json:"prompt,omitempty"`
	DynamicVariablePlaceholders map[string]string `json:"dynamic_variable_placeholders,omitempty"`
}

type ConversationalConfig struct {
	ASR             *ASRConfig         `json:"asr,omitempty"`
	Turn            *TurnConfig        `json:"turn,omitempty"`
	TTS             *AgentTTSConfig    `json:"tts,omitempty"`
	Conversation    *AgentConversation `json:"conversation,omitempty"`
	LanguagePresets json.RawMessage    `json:"language_presets,omitempty"`
	VAD             json.RawMessage    `json:"vad,omitempty"`
	Agent           *AgentConfig       `json:"agent,omitempty"`
}

type Agent struct {
	AgentID             string               `json:"agent_id"`
	Name                string               `json:"name"`
	ConversationConfig  ConversationalConfig `json:"conversation_config"`
	Metadata            AgentMetadata        `json:"metadata"`
	PlatformSettings    json.RawMessage      `json:"platform_settings,omitempty"`
	PhoneNumbers        json.RawMessage      `json:"phone_numbers,omitempty"`
	WhatsAppAccounts    json.RawMessage      `json:"whatsapp_accounts,omitempty"`
	Workflow            json.RawMessage      `json:"workflow,omitempty"`
	AccessInfo          *ResourceAccessInfo  `json:"access_info,omitempty"`
	Tags                []string             `json:"tags,omitempty"`
	VersionID           string               `json:"version_id,omitempty"`
	BranchID            string               `json:"branch_id,omitempty"`
	MainBranchID        string               `json:"main_branch_id,omitempty"`
	Procedures          json.RawMessage      `json:"procedures,omitempty"`
	DefaultHoldAudioURL string               `json:"default_hold_audio_url,omitempty"`
}

type CreateAgentResult struct {
	AgentID string `json:"agent_id"`
}

type ConversationLinkToken struct {
	AgentID                string `json:"agent_id"`
	ConversationToken      string `json:"conversation_token"`
	Purpose                string `json:"purpose"`
	ExpirationTimeUnixSecs *int64 `json:"expiration_time_unix_secs"`
	ConversationID         string `json:"conversation_id,omitempty"`
	TokenRequesterUserID   string `json:"token_requester_user_id,omitempty"`
}

type AgentLink struct {
	AgentID string                 `json:"agent_id"`
	Token   *ConversationLinkToken `json:"token,omitempty"`
}
