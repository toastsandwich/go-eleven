# Agents

Create, list, create-or-update, update, duplicate, and delete Eleven Agents with `/v1/convai/agents`.

This is agent configuration only. For past calls and session tokens, see [conversations](conversations.md). Workspace tools are in [tools](tools.md). Documents are in [knowledge-base](knowledge-base.md). Live sessions are in [live-conversation](live-conversation.md).

[Back to docs](README.md)

## Create agent

`CreateAgent` takes a `ConversationalConfig` and returns the new `agent_id`. Name and tags are options.

```go
created, err := client.CreateAgent(
	eleven.ConversationalConfig{
		Agent: &eleven.AgentConfig{
			FirstMessage: "Hi, how can I help?",
			Language:     "en",
			Prompt: &eleven.AgentPrompt{
				Prompt: "You are a concise support agent.",
				LLM:    "gemini-2.0-flash",
			},
		},
		TTS: &eleven.AgentTTSConfig{
			VoiceID: voiceId,
			ModelID: "eleven_turbo_v2_5",
		},
	},
	eleven.WithAgentName("Support Concierge"),
	eleven.WithCreateAgentTags("support"),
)
if err != nil {
	log.Fatal(err)
}
fmt.Println(created.AgentID)
```

`platform_settings` and `workflow` are pass-through JSON (`WithPlatformSettings`, `WithAgentWorkflow`).

## Create or update agent

`CreateOrUpdateAgent` looks up the agent by exact name. Missing means create; one match means patch `conversation_config` plus any tags, platform settings, and workflow you pass. `WithAgentName` is required. A substring list hit with a different name is ignored.

```go
result, err := client.CreateOrUpdateAgent(
	eleven.ConversationalConfig{
		Agent: &eleven.AgentConfig{
			FirstMessage: "Hi, how can I help?",
			Language:     "en",
			Prompt: &eleven.AgentPrompt{
				Prompt: "You are a concise support agent.",
				LLM:    "gemini-2.0-flash",
			},
		},
		TTS: &eleven.AgentTTSConfig{
			VoiceID: voiceId,
			ModelID: "eleven_turbo_v2_5",
		},
	},
	eleven.WithAgentName("Support Concierge"),
	eleven.WithCreateAgentTags("support"),
)
if err != nil {
	log.Fatal(err)
}
fmt.Println(result.AgentID, result.Created)
```

Two agents with the same name is an error. Concurrent calls can still insert duplicates because the API does not unique-index names.

## List agents

```go
result, err := client.ListAgents(
	eleven.WithAgentSearch("Support"),
	eleven.WithAgentPageSize(30),
	eleven.WithAgentSort(eleven.AgentSortName),
	eleven.WithAgentSortDirection(eleven.SortAsc),
)
if err != nil {
	log.Fatal(err)
}

for agent := range result.Agents.All() {
	fmt.Println(agent.Name, agent.AgentID)
}
```

`Agents` also has `AllWith` for extra local filters after the API response.

Paginate with `has_more` and `next_cursor`.

```go
if result.HasMore {
	result, err = client.ListAgents(
		eleven.WithAgentPageSize(30),
		eleven.WithAgentCursor(result.NextCursor),
	)
}
```

## Get, update, duplicate, link, delete

```go
agent, err := client.GetAgent(agentID)

agent, err = client.UpdateAgent(agentID, eleven.WithUpdateAgentName("Support v2"))

copy, err := client.DuplicateAgent(agentID, eleven.WithDuplicateAgentName("Support copy"))

link, err := client.GetAgentLink(agentID)

err = client.DeleteAgent(copy.AgentID, eleven.WithForceDeleteAgent())
err = client.DeleteAgent(agentID)
```

A name-only update does not send `conversation_config`. Use `WithUpdateConversationConfig` to patch prompt, voice, or turn settings.

`ConversationalConfig` also carries `Turn.SoftTimeoutConfig`, `Agent.DynamicVariablePlaceholders`, and `Prompt.MCPServerIDs` when you PATCH those ElevenLabs fields.

## Options

### List

| Option | Query | Notes |
| --- | --- | --- |
| `WithAgentPageSize` | `page_size` | Max 100, default 30 |
| `WithAgentSearch` | `search` | Match agent name |
| `WithAgentArchived` | `archived` | Archived agents only |
| `WithCreatedByUserID` | `created_by_user_id` | `@me` is the authenticated user |
| `WithAgentTags` | `tags` | Repeated; match any |
| `WithAgentSort` | `sort_by` | `name`, `created_at`, `updated_at`, `call_count_7d` |
| `WithAgentSortDirection` | `sort_direction` | `SortAsc` / `SortDesc` |
| `WithAgentCursor` | `cursor` | Next page from `NextCursor` |

### Get

| Option | Query | Notes |
| --- | --- | --- |
| `WithAgentVersionID` | `version_id` | |
| `WithAgentBranchID` | `branch_id` | |
| `WithoutProcedureGeneratedWorkflow` | `include_procedure_generated_workflow_elements=false` | Authored workflow only |

### Create

| Option | Body | Notes |
| --- | --- | --- |
| `WithAgentName` | `name` | |
| `WithCreateAgentTags` | `tags` | |
| `WithPlatformSettings` | `platform_settings` | Raw JSON |
| `WithAgentWorkflow` | `workflow` | Raw JSON |

### Create or update

Same options as Create. `WithAgentName` is required and is the lookup key.

### Update

| Option | Where | Notes |
| --- | --- | --- |
| `WithUpdateAgentName` | body | |
| `WithUpdateAgentTags` | body | Replaces tags |
| `WithUpdateConversationConfig` | body | |
| `WithUpdatePlatformSettings` | body | Raw JSON |
| `WithUpdateAgentWorkflow` | body | Raw JSON |
| `WithVersionDescription` | body | Version note when publishing |
| `WithUpdateAgentBranchID` | query | |
| `WithoutUpdateProcedureGeneratedWorkflow` | query | Authored workflow only |

### Delete / duplicate

| Option | Notes |
| --- | --- |
| `WithForceDeleteAgent` | Delete even if used as a subagent |
| `WithDuplicateAgentName` | Name for the copy |

## Context

```go
result, err := client.ListAgentsWithContext(ctx)
created, err := client.CreateAgentWithContext(ctx, config, eleven.WithAgentName("Support"))
result, err := client.CreateOrUpdateAgentWithContext(ctx, config, eleven.WithAgentName("Support"))
agent, err := client.GetAgentWithContext(ctx, agentID)
agent, err = client.UpdateAgentWithContext(ctx, agentID, eleven.WithUpdateAgentName("Support v2"))
err = client.DeleteAgentWithContext(ctx, agentID)
```
