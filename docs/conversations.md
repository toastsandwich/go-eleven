# Conversations

List past Eleven Agent calls, fetch transcripts and audio, and mint a session for a client with `/v1/convai/conversations` and `/v1/convai/conversation`.

Open a live WebSocket from this process with [live conversation](live-conversation.md). Give a signed URL or WebRTC token to a browser or mobile client SDK.

Post-call webhooks use a different payload. Verify those with `ConstructWebhookEvent` and `PostCallTranscriptionData`. See [webhooks](webhooks.md).

[Back to docs](README.md)

## List conversations

```go
result, err := client.ListConversations(
	eleven.WithConversationAgentID(agentID),
	eleven.WithConversationPageSize(30),
	eleven.WithCallSuccessful(eleven.CallSuccessfulSuccess),
)
if err != nil {
	log.Fatal(err)
}

for conv := range result.Conversations.All() {
	fmt.Println(conv.ConversationID, conv.Status, conv.CallDurationSecs)
}
```

`Conversations` also has `AllWith` for extra local filters after the API response.

Paginate with `has_more` and `next_cursor`.

```go
if result.HasMore {
	result, err = client.ListConversations(
		eleven.WithConversationAgentID(agentID),
		eleven.WithConversationPageSize(30),
		eleven.WithConversationCursor(result.NextCursor),
	)
}
```

## Get details, summary, and audio

```go
details, err := client.GetConversation(conversationID)
summary, err := client.GetConversationSummary(conversationID, eleven.WithMaxMessages(20))
audio, err := client.GetConversationAudio(conversationID)
```

`GetConversation` returns the full transcript, including tool calls as JSON. `GetConversationSummary` returns title, gist, and plain chat messages when the call is short. `GetConversationAudio` returns the recording as bytes. The HTTP timeout is disabled so longer recordings can finish. Cancel with context.

## Feedback and delete

```go
err = client.SendConversationFeedback(conversationID, eleven.ConversationFeedbackLike)
err = client.DeleteConversation(conversationID)
```

## Start a client session

Mint a signed URL for an agent that requires authorization, or a WebRTC token for realtime audio. Pass the signed URL to `StartConversation` with `WithSessionURL`, or let the client fetch it with `WithSessionSignedURL`. See [live conversation](live-conversation.md).

```go
signed, err := client.GetSignedURL(agentID, eleven.WithSignedURLConversationID())
token, err := client.GetConversationToken(agentID, eleven.WithParticipantName("web-user"))
```

For OpenAI-compatible Realtime browser sessions:

```go
secret, err := client.CreateRealtimeClientSecret(eleven.WithRealtimeExpiresAfter(600))
```

Session config for the Realtime secret is pass-through JSON (`WithRealtimeSession`).

## Options

### List

| Option | Query | Notes |
| --- | --- | --- |
| `WithConversationCursor` | `cursor` | Next page from `NextCursor` |
| `WithConversationAgentID` | `agent_id` | Restrict to one agent |
| `WithConversationPageSize` | `page_size` | Max 100, default 30 |
| `WithConversationSortDirection` | `sort_direction` | `SortAsc` / `SortDesc` |
| `WithCallSuccessful` | `call_successful` | `success`, `failure`, `unknown`, `error` |
| `WithCallStartAfterUnix` | `call_start_after_unix` | |
| `WithCallStartBeforeUnix` | `call_start_before_unix` | |
| `WithConversationUserID` | `user_id` | |
| `WithConversationBranchID` | `branch_id` | |
| `WithConversationVersionID` | `version_id` | |
| `WithSummaryMode` | `summary_mode` | `exclude` (default) or `include` |

### Summary

| Option | Query | Notes |
| --- | --- | --- |
| `WithMaxMessages` | `max_messages` | Default 40. Longer calls omit messages |

### Signed URL

| Option | Query | Notes |
| --- | --- | --- |
| `WithSignedURLConversationID` | `include_conversation_id=true` | Token cannot be reused |
| `WithSignedURLBranchID` | `branch_id` | |
| `WithSignedURLVersionID` | `version_id` | |
| `WithSignedURLEnvironment` | `environment` | Default `production` |

### WebRTC token

| Option | Query | Notes |
| --- | --- | --- |
| `WithParticipantName` | `participant_name` | |
| `WithTokenBranchID` | `branch_id` | |
| `WithTokenVersionID` | `version_id` | |
| `WithTokenEnvironment` | `environment` | |

### Realtime secret

| Option | Body | Notes |
| --- | --- | --- |
| `WithRealtimeSession` | `session` | Raw JSON |
| `WithRealtimeExpiresAfter` | `expires_after.seconds` | Anchored at `created_at` |

## Context

```go
list, err := client.ListConversationsWithContext(ctx, eleven.WithConversationAgentID(agentID))
details, err := client.GetConversationWithContext(ctx, conversationID)
signed, err := client.GetSignedURLWithContext(ctx, agentID)
```
