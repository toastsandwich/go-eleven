package eleven

import (
	"encoding/base64"
	"encoding/json"
)

type SessionEventType string

const (
	SessionEventInitiation              SessionEventType = "conversation_initiation_metadata"
	SessionEventPing                    SessionEventType = "ping"
	SessionEventAudio                   SessionEventType = "audio"
	SessionEventUserTranscript          SessionEventType = "user_transcript"
	SessionEventAgentResponse           SessionEventType = "agent_response"
	SessionEventAgentResponseCorrection SessionEventType = "agent_response_correction"
	SessionEventAgentChatResponsePart   SessionEventType = "agent_chat_response_part"
	SessionEventAgentResponseComplete   SessionEventType = "agent_response_complete"
	SessionEventInterruption            SessionEventType = "interruption"
	SessionEventClientToolCall          SessionEventType = "client_tool_call"
	SessionEventAgentToolResponse       SessionEventType = "agent_tool_response"
	SessionEventVADScore                SessionEventType = "vad_score"
	SessionEventQueueStatus             SessionEventType = "queue_status"
	SessionEventClientError             SessionEventType = "client_error"
	SessionEventContextUsage            SessionEventType = "context_usage"
	SessionEventGuardrailTriggered      SessionEventType = "guardrail_triggered"
	SessionEventUnknown                 SessionEventType = "unknown"
)

func (t SessionEventType) String() string { return string(t) }

type SessionChatPartType string

const (
	SessionChatPartStart SessionChatPartType = "start"
	SessionChatPartDelta SessionChatPartType = "delta"
	SessionChatPartStop  SessionChatPartType = "stop"
)

func (t SessionChatPartType) String() string { return string(t) }

type SessionQueueStatusValue string

const (
	SessionQueueWaiting  SessionQueueStatusValue = "waiting"
	SessionQueueAdmitted SessionQueueStatusValue = "admitted"
	SessionQueueTimedOut SessionQueueStatusValue = "timed_out"
)

func (s SessionQueueStatusValue) String() string { return string(s) }

type SessionInitiation struct {
	ConversationID         string `json:"conversation_id"`
	AgentOutputAudioFormat string `json:"agent_output_audio_format"`
	UserInputAudioFormat   string `json:"user_input_audio_format"`
}

type SessionPing struct {
	EventID int `json:"event_id"`
	PingMS  int `json:"ping_ms,omitempty"`
}

type SessionAudio struct {
	PCM       []byte          `json:"-"`
	EventID   int             `json:"event_id"`
	Alignment json.RawMessage `json:"alignment,omitempty"`
}

type SessionTranscript struct {
	Text    string `json:"user_transcript"`
	EventID int    `json:"event_id,omitempty"`
}

type SessionAgentResponse struct {
	Text       string `json:"agent_response"`
	EventID    int    `json:"event_id,omitempty"`
	ResponseID string `json:"response_id,omitempty"`
}

type SessionCorrection struct {
	Original  string `json:"original_agent_response"`
	Corrected string `json:"corrected_agent_response"`
}

type SessionChatPart struct {
	Type       SessionChatPartType `json:"type"`
	Text       string              `json:"text,omitempty"`
	EventID    int                 `json:"event_id,omitempty"`
	ResponseID string              `json:"response_id,omitempty"`
}

type SessionResponseComplete struct {
	EventID int `json:"event_id"`
}

type SessionInterruption struct {
	Reason string `json:"reason,omitempty"`
}

type SessionClientToolCall struct {
	ToolName   string          `json:"tool_name"`
	ToolCallID string          `json:"tool_call_id"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

type SessionAgentToolResponse struct {
	ToolName   string `json:"tool_name"`
	ToolCallID string `json:"tool_call_id"`
	ToolType   string `json:"tool_type,omitempty"`
	IsError    bool   `json:"is_error,omitempty"`
}

type SessionVAD struct {
	Score float64 `json:"vad_score"`
}

type SessionQueueStatus struct {
	Status SessionQueueStatusValue `json:"status"`
}

type SessionClientError struct {
	Code      int    `json:"code"`
	ErrorName string `json:"error_name"`
	Message   string `json:"message,omitempty"`
}

type SessionContextUsage struct {
	EventID            int    `json:"event_id"`
	Model              string `json:"model"`
	ContextTokens      int    `json:"context_tokens"`
	ContextLimitTokens int    `json:"context_limit_tokens"`
}

type SessionGuardrail struct {
	Name string `json:"guardrail_name"`
}

type SessionEvent struct {
	Type SessionEventType
	Raw  json.RawMessage

	Initiation   *SessionInitiation
	Ping         *SessionPing
	Audio        *SessionAudio
	Transcript   *SessionTranscript
	Response     *SessionAgentResponse
	Correction   *SessionCorrection
	ChatPart     *SessionChatPart
	Complete     *SessionResponseComplete
	Interruption *SessionInterruption
	ClientTool   *SessionClientToolCall
	AgentTool    *SessionAgentToolResponse
	VAD          *SessionVAD
	Queue        *SessionQueueStatus
	Error        *SessionClientError
	ContextUsage *SessionContextUsage
	Guardrail    *SessionGuardrail
}

type ConversationSession interface {
	Events() <-chan SessionEvent
	ConversationID() string
	SendUserMessage(text string) error
	SendAudio(pcm []byte) error
	SendContextualUpdate(text string) error
	SendUserActivity() error
	SendClientToolResult(toolCallID string, result any, isError bool) error
	SendFeedback(eventID int, feedback ConversationFeedback) error
	Close() error
	Err() error
}

func decodeSessionEvent(body []byte) (SessionEvent, error) {
	event := SessionEvent{
		Raw:  json.RawMessage(append([]byte(nil), body...)),
		Type: SessionEventUnknown,
	}

	var head struct {
		Type SessionEventType `json:"type"`
	}
	if err := json.Unmarshal(body, &head); err != nil {
		return event, err
	}
	if head.Type == "" {
		return event, nil
	}

	switch head.Type {
	case SessionEventInitiation:
		var frame struct {
			Event SessionInitiation `json:"conversation_initiation_metadata_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Initiation = &frame.Event
	case SessionEventPing:
		var frame struct {
			Event SessionPing `json:"ping_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Ping = &frame.Event
	case SessionEventAudio:
		var frame struct {
			Event struct {
				AudioBase64 string          `json:"audio_base_64"`
				EventID     int             `json:"event_id"`
				Alignment   json.RawMessage `json:"alignment"`
			} `json:"audio_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		pcm, err := base64.StdEncoding.DecodeString(frame.Event.AudioBase64)
		if err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Audio = &SessionAudio{
			PCM:       pcm,
			EventID:   frame.Event.EventID,
			Alignment: frame.Event.Alignment,
		}
	case SessionEventUserTranscript:
		var frame struct {
			Event SessionTranscript `json:"user_transcription_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Transcript = &frame.Event
	case SessionEventAgentResponse:
		var frame struct {
			Event SessionAgentResponse `json:"agent_response_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Response = &frame.Event
	case SessionEventAgentResponseCorrection:
		var frame struct {
			Event SessionCorrection `json:"agent_response_correction_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Correction = &frame.Event
	case SessionEventAgentChatResponsePart:
		var frame struct {
			Event SessionChatPart `json:"text_response_part"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.ChatPart = &frame.Event
	case SessionEventAgentResponseComplete:
		var frame struct {
			Event SessionResponseComplete `json:"agent_response_complete_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Complete = &frame.Event
	case SessionEventInterruption:
		var frame struct {
			Event SessionInterruption `json:"interruption_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Interruption = &frame.Event
	case SessionEventClientToolCall:
		var frame struct {
			Event SessionClientToolCall `json:"client_tool_call"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.ClientTool = &frame.Event
	case SessionEventAgentToolResponse:
		var frame struct {
			Event SessionAgentToolResponse `json:"agent_tool_response"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.AgentTool = &frame.Event
	case SessionEventVADScore:
		var frame struct {
			Event SessionVAD `json:"vad_score_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.VAD = &frame.Event
	case SessionEventQueueStatus:
		var frame struct {
			Event SessionQueueStatus `json:"queue_status_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Queue = &frame.Event
	case SessionEventClientError:
		var frame struct {
			Event SessionClientError `json:"error_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Error = &frame.Event
	case SessionEventContextUsage:
		var frame struct {
			Event SessionContextUsage `json:"context_usage_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.ContextUsage = &frame.Event
	case SessionEventGuardrailTriggered:
		var frame struct {
			Event SessionGuardrail `json:"guardrail_triggered_event"`
		}
		if err := json.Unmarshal(body, &frame); err != nil {
			return event, err
		}
		event.Type = head.Type
		event.Guardrail = &frame.Event
	default:
		event.Type = SessionEventUnknown
	}
	return event, nil
}
