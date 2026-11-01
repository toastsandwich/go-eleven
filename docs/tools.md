# Tools

Create, list, update, and delete workspace tools with `/v1/convai/tools`. Attach a tool to an agent with `AgentPrompt.ToolIDs`. See [agents](agents.md).

[Back to docs](README.md)

## Create tool

Client tools run in the conversation client. Webhook tools call an HTTPS endpoint.

```go
clientTool, err := client.CreateTool(eleven.ToolConfig{
	Type:        eleven.ToolTypeClient,
	Name:        "show_order_status",
	Description: "Show the user's order status in the client UI.",
})
if err != nil {
	log.Fatal(err)
}

webhookTool, err := client.CreateTool(eleven.ToolConfig{
	Type:        eleven.ToolTypeWebhook,
	Name:        "lookup_order",
	Description: "Fetch order status from the backend.",
	APISchema: &eleven.WebhookAPISchema{
		URL:    "https://example.com/orders",
		Method: "GET",
	},
})
```

`parameters` (client JSON schema), `params` (system tools), `value` (MCP), `assignments`, and `dynamic_variables` are pass-through JSON. Mock responses use `WithToolResponseMocks`.

Create always inserts a new tool. Persist `id` if you need the same tool after a restart.

## List tools

```go
result, err := client.ListTools(
	eleven.WithToolSearch("show_"),
	eleven.WithToolPageSize(30),
	eleven.WithToolTypes(eleven.ToolTypeClient, eleven.ToolTypeWebhook),
	eleven.WithToolSort(eleven.ToolSortName),
	eleven.WithToolSortDirection(eleven.SortAsc),
)
if err != nil {
	log.Fatal(err)
}

for tool := range result.Tools.All() {
	fmt.Println(tool.ID, tool.ToolConfig.Name, tool.ToolConfig.Type)
}
```

`Tools` also has `AllWith` for extra local filters after the API response.

Paginate with `has_more` and `next_cursor`.

```go
if result.HasMore {
	result, err = client.ListTools(
		eleven.WithToolPageSize(30),
		eleven.WithToolCursor(result.NextCursor),
	)
}
```

## Get, update, dependents, delete

```go
tool, err := client.GetTool(toolID)
tool, err = client.UpdateTool(toolID, eleven.ToolConfig{
	Type:        eleven.ToolTypeClient,
	Name:        "show_order_status",
	Description: "Updated description.",
})

dependents, err := client.GetToolDependentAgents(toolID)
err = client.DeleteTool(toolID, eleven.WithForceDeleteTool())
err = client.DeleteTool(toolID)
```

Update sends a full `tool_config`. Delete fails while agents still reference the tool unless you pass `WithForceDeleteTool`.

## Attach to an agent

```go
_, err = client.UpdateAgent(agentID, eleven.WithUpdateConversationConfig(eleven.ConversationalConfig{
	Agent: &eleven.AgentConfig{
		Prompt: &eleven.AgentPrompt{
			ToolIDs: []string{clientTool.ID, webhookTool.ID},
		},
	},
}))
```

## Options

### List

| Option | Query | Notes |
| --- | --- | --- |
| `WithToolSearch` | `search` | Name prefix |
| `WithToolPageSize` | `page_size` | Max 100, default 30 |
| `WithToolCreatedByUserID` | `created_by_user_id` | `@me` is the authenticated user |
| `WithToolTypes` | `types` | Repeated; `webhook`, `client`, `api_integration_webhook` |
| `WithToolSort` | `sort_by` | `name`, `created_at` |
| `WithToolSortDirection` | `sort_direction` | `SortAsc` / `SortDesc` |
| `WithToolUsedByAgentID` | `used_by_agent_id` | Tools on that agent's main branch |
| `WithToolCursor` | `cursor` | Next page from `NextCursor` |

### Get

| Option | Query | Notes |
| --- | --- | --- |
| `WithToolEnvironment` | `environment` | MCP env vars; default `production` |

### Create / update

| Option | Body | Notes |
| --- | --- | --- |
| `WithToolResponseMocks` | `response_mocks` | Raw JSON |

### Delete / dependents

| Option | Notes |
| --- | --- |
| `WithForceDeleteTool` | Delete even if agents use the tool |
| `WithToolDependentCursor` | Next page from `NextCursor` |
| `WithToolDependentPageSize` | Max 100, default 30 |

## Context

```go
list, err := client.ListToolsWithContext(ctx, eleven.WithToolSearch("show_"))
tool, err := client.CreateToolWithContext(ctx, config)
tool, err = client.GetToolWithContext(ctx, toolID)
err = client.DeleteToolWithContext(ctx, toolID)
```
