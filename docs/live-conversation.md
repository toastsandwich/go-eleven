# Live conversation

Talk to an Eleven Agent over `wss://api.elevenlabs.io/v1/convai/conversation`.

This is a protocol client: send PCM or text, receive typed events. It does not capture a microphone, play speakers, or open WebRTC. For past calls and session tokens, see [conversations](conversations.md).

[Back to docs](README.md)

## Start a session

`StartConversation` dials the socket, writes `conversation_initiation_client_data`, and returns a `ConversationSession`. The read loop runs until `Close` or context cancel.

```go
session, err := client.StartConversation(agentID)
if err != nil {
	log.Fatal(err)
}
defer session.Close()

go func() {
	for event := range session.Events() {
		switch event.Type {
		case eleven.SessionEventInitiation:
			fmt.Println(event.Initiation.ConversationID)
		case eleven.SessionEventAgentResponse:
			fmt.Println(event.Response.Text)
		case eleven.SessionEventAudio:
			_ = event.Audio.PCM
		case eleven.SessionEventClientToolCall:
			_ = session.SendClientToolResult(event.ClientTool.ToolCallID, "ok", false)
		}
	}
}()

if err := session.SendUserMessage("What are your hours?"); err != nil {
	log.Fatal(err)
}
```

Drain `Events()`. The channel is buffered (default 32). `ConversationID()` is empty until `conversation_initiation_metadata`. After the channel closes, `Err()` is the read or close error. A clean close is nil.

The HTTP client timeout is disabled for the socket. Cancel with context.

## Signed URL

Private agents need a signed URL. `WithSessionURL` dials that URL and skips `xi-api-key`. Empty `agentID` is allowed when the URL is set.

```go
signed, err := client.GetSignedURL(agentID)
session, err := client.StartConversation("", eleven.WithSessionURL(signed.SignedURL))
```

`WithSessionSignedURL` fetches the URL with `GetSignedURL` first. `WithSessionURL` wins if both are set.

```go
session, err := client.StartConversation(agentID, eleven.WithSessionSignedURL())
```

Public agents can dial with `agent_id` and the API key.

## Send

| Method | Frame |
| --- | --- |
| (auto, on connect) | `conversation_initiation_client_data` |
| `SendUserMessage` | `user_message` |
| `SendAudio` | `user_audio_chunk` (PCM, base64 encoded internally) |
| `SendContextualUpdate` | `contextual_update` |
| `SendUserActivity` | `user_activity` |
| `SendClientToolResult` | `client_tool_result` |
| `SendFeedback` | `feedback` (`like` / `dislike`) |
| (auto ping) | `pong` |

`SendAudio` takes raw PCM. The agent’s input format is on `SessionInitiation.UserInputAudioFormat` (`pcm_16000`, `pcm_44100`, …). Server `ping` is answered automatically; the event is still delivered.

Empty checks: `"text is empty"`, `"audio is empty"`, `"toolCallID is empty"`, `"eventID is empty"`. Empty `agentID` is `"agentID is empty"` unless `WithSessionURL` is set.

## Events

| Type | Payload |
| --- | --- |
| `conversation_initiation_metadata` | `Initiation` |
| `ping` | `Ping` |
| `audio` | `Audio` (decoded PCM) |
| `user_transcript` | `Transcript` |
| `agent_response` | `Response` |
| `agent_response_correction` | `Correction` |
| `agent_chat_response_part` | `ChatPart` |
| `agent_response_complete` | `Complete` |
| `interruption` | `Interruption` |
| `client_tool_call` | `ClientTool` |
| `agent_tool_response` | `AgentTool` |
| `vad_score` | `VAD` |
| `queue_status` | `Queue` |
| `client_error` | `Error` |
| `context_usage` | `ContextUsage` |
| `guardrail_triggered` | `Guardrail` |
| `unknown` | parse `Raw` |

`mcp_tool_call` and `agent_response_metadata` land as `SessionEventUnknown` with `Raw` set. Every event keeps the original frame on `Raw`.

## Options

| Option | Where | Notes |
| --- | --- | --- |
| `WithSessionURL` | dial | Signed URL from `GetSignedURL`; skip `xi-api-key` |
| `WithSessionSignedURL` | dial | Call `GetSignedURL`, then dial |
| `WithSessionBranchID` | query | |
| `WithSessionVersionID` | query | |
| `WithSessionUserID` | initiation `user_id` | |
| `WithSessionConfigOverride` | initiation | Raw JSON |
| `WithSessionDynamicVariables` | initiation | Raw JSON |
| `WithSessionExtraBody` | initiation `custom_llm_extra_body` | Raw JSON |
| `WithSessionEventBuffer` | channel size | Default 32 |

## Context

```go
session, err := client.StartConversationWithContext(ctx, agentID, eleven.WithSessionUserID("user_1"))
```
