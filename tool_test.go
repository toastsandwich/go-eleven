package eleven

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestToolEnums(t *testing.T) {
	if ToolTypeWebhook.String() != "webhook" {
		t.Fatalf("type = %s", ToolTypeWebhook)
	}
	if ToolTypeAPIIntegrationWebhook.String() != "api_integration_webhook" {
		t.Fatalf("type = %s", ToolTypeAPIIntegrationWebhook)
	}
	if ToolSortCreatedAt.String() != "created_at" {
		t.Fatalf("sort = %s", ToolSortCreatedAt)
	}
	if ToolExecutionPostToolSpeech.String() != "post_tool_speech" {
		t.Fatalf("execution = %s", ToolExecutionPostToolSpeech)
	}
	if ToolInterruptionDisableDuringToolAndTurn.String() != "disable_during_tool_and_turn" {
		t.Fatalf("interruption = %s", ToolInterruptionDisableDuringToolAndTurn)
	}
	if ToolPreSpeechForce.String() != "force" {
		t.Fatalf("pre speech = %s", ToolPreSpeechForce)
	}
	if ToolCallSoundElevator2.String() != "elevator2" {
		t.Fatalf("sound = %s", ToolCallSoundElevator2)
	}
	if ToolCallSoundAlways.String() != "always" {
		t.Fatalf("behavior = %s", ToolCallSoundAlways)
	}
	if ToolErrorHandlingPassthrough.String() != "passthrough" {
		t.Fatalf("error handling = %s", ToolErrorHandlingPassthrough)
	}
	if DependentAgentAvailable.String() != "available" {
		t.Fatalf("dependent = %s", DependentAgentAvailable)
	}
}

func TestListToolsQuery(t *testing.T) {
	var opt listToolsOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithToolSearch("show_")(&opt)
	WithToolPageSize(20)(&opt)
	WithToolCreatedByUserID("@me")(&opt)
	WithToolTypes(ToolTypeWebhook, ToolTypeClient)(&opt)
	WithToolSort(ToolSortName)(&opt)
	WithToolSortDirection(SortAsc)(&opt)
	WithToolUsedByAgentID("agent_1")(&opt)
	WithToolCursor("cur_1")(&opt)

	q := opt.query()
	if q.Get("search") != "show_" {
		t.Fatalf("search = %q", q.Get("search"))
	}
	if q.Get("page_size") != "20" {
		t.Fatalf("page_size = %q", q.Get("page_size"))
	}
	if q.Get("created_by_user_id") != "@me" {
		t.Fatalf("created_by_user_id = %q", q.Get("created_by_user_id"))
	}
	if got := q["types"]; len(got) != 2 || got[0] != "webhook" || got[1] != "client" {
		t.Fatalf("types = %v", got)
	}
	if q.Get("sort_by") != "name" {
		t.Fatalf("sort_by = %q", q.Get("sort_by"))
	}
	if q.Get("sort_direction") != "asc" {
		t.Fatalf("sort_direction = %q", q.Get("sort_direction"))
	}
	if q.Get("used_by_agent_id") != "agent_1" {
		t.Fatalf("used_by_agent_id = %q", q.Get("used_by_agent_id"))
	}
	if q.Get("cursor") != "cur_1" {
		t.Fatalf("cursor = %q", q.Get("cursor"))
	}
}

func TestGetToolQuery(t *testing.T) {
	var opt getToolOption
	WithToolEnvironment("staging")(&opt)
	if opt.query().Get("environment") != "staging" {
		t.Fatalf("query = %v", opt.query())
	}
}

func TestCreateToolRequestBody(t *testing.T) {
	var opt createToolOption
	got, err := json.Marshal(opt.body(ToolConfig{
		Type:        ToolTypeClient,
		Name:        "show_order_status",
		Description: "Show order status.",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"tool_config":{"type":"client","name":"show_order_status","description":"Show order status."}}` {
		t.Fatalf("minimal body = %s", got)
	}

	WithToolResponseMocks(json.RawMessage(`[{"mock_result":"ok"}]`))(&opt)
	got, err = json.Marshal(opt.body(ToolConfig{Type: ToolTypeClient, Name: "show_order_status"}))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	mocks := payload["response_mocks"].([]any)
	if len(mocks) != 1 {
		t.Fatalf("response_mocks = %v", payload["response_mocks"])
	}
}

func TestDeleteToolQuery(t *testing.T) {
	var opt deleteToolOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}
	WithForceDeleteTool()(&opt)
	if opt.query().Get("force") != "true" {
		t.Fatalf("query = %v", opt.query())
	}
}

func TestToolDependentQuery(t *testing.T) {
	var opt toolDependentOption
	WithToolDependentCursor("cur_1")(&opt)
	WithToolDependentPageSize(10)(&opt)
	q := opt.query()
	if q.Get("cursor") != "cur_1" || q.Get("page_size") != "10" {
		t.Fatalf("query = %v", q)
	}
}

func TestToolValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.GetTool(""); err == nil || err.Error() != "toolID is empty" {
		t.Fatalf("get err = %v", err)
	}
	if _, err := client.UpdateTool("", ToolConfig{}); err == nil || err.Error() != "toolID is empty" {
		t.Fatalf("update err = %v", err)
	}
	if err := client.DeleteTool(""); err == nil || err.Error() != "toolID is empty" {
		t.Fatalf("delete err = %v", err)
	}
	if _, err := client.GetToolDependentAgents(""); err == nil || err.Error() != "toolID is empty" {
		t.Fatalf("dependent err = %v", err)
	}
}

func TestListTools(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("search") != "show_" || r.URL.Query().Get("page_size") != "20" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"tools":[{"id":"tool_1","tool_config":{"type":"client","name":"show_order_status","description":"Show status."},"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"usage_stats":{"avg_latency_secs":0.2,"total_calls":3}}],"has_more":true,"next_cursor":"cur_2"}`))
	})

	got, err := client.ListTools(WithToolSearch("show_"), WithToolPageSize(20))
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasMore || got.NextCursor != "cur_2" || len(got.Tools) != 1 {
		t.Fatalf("list = %+v", got)
	}
	if got.Tools[0].ID != "tool_1" || got.Tools[0].ToolConfig.Name != "show_order_status" {
		t.Fatalf("tool = %+v", got.Tools[0])
	}
}

func TestCreateTool(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		var payload createToolRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.ToolConfig.Type != ToolTypeClient || payload.ToolConfig.Name != "show_order_status" {
			t.Fatalf("config = %+v", payload.ToolConfig)
		}
		w.Write([]byte(`{"id":"tool_1","tool_config":{"type":"client","name":"show_order_status","description":"Show the user's order status in the client UI."},"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"usage_stats":{"avg_latency_secs":0,"total_calls":0}}`))
	})

	got, err := client.CreateTool(ToolConfig{
		Type:        ToolTypeClient,
		Name:        "show_order_status",
		Description: "Show the user's order status in the client UI.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "tool_1" || got.ToolConfig.Type != ToolTypeClient {
		t.Fatalf("result = %+v", got)
	}
}

func TestCreateWebhookTool(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		var payload createToolRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.ToolConfig.Type != ToolTypeWebhook || payload.ToolConfig.APISchema == nil || payload.ToolConfig.APISchema.URL != "https://example.com/orders" {
			t.Fatalf("config = %+v", payload.ToolConfig)
		}
		w.Write([]byte(`{"id":"tool_2","tool_config":{"type":"webhook","name":"lookup_order","description":"Fetch order status.","api_schema":{"url":"https://example.com/orders","method":"GET"}},"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"usage_stats":{"avg_latency_secs":0,"total_calls":0}}`))
	})

	got, err := client.CreateTool(ToolConfig{
		Type:        ToolTypeWebhook,
		Name:        "lookup_order",
		Description: "Fetch order status.",
		APISchema: &WebhookAPISchema{
			URL:    "https://example.com/orders",
			Method: "GET",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "tool_2" || got.ToolConfig.APISchema == nil || got.ToolConfig.APISchema.URL != "https://example.com/orders" {
		t.Fatalf("result = %+v", got)
	}
}

func TestGetTool(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools/tool_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("environment") != "staging" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"id":"tool_1","tool_config":{"type":"client","name":"show_order_status"},"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"usage_stats":{"avg_latency_secs":0.1,"total_calls":2}}`))
	})

	got, err := client.GetTool("tool_1", WithToolEnvironment("staging"))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "tool_1" || got.UsageStats.TotalCalls != 2 {
		t.Fatalf("tool = %+v", got)
	}
}

func TestUpdateTool(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools/tool_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload createToolRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.ToolConfig.Description != "Updated description." {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"id":"tool_1","tool_config":{"type":"client","name":"show_order_status","description":"Updated description."},"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"usage_stats":{"avg_latency_secs":0,"total_calls":0}}`))
	})

	got, err := client.UpdateTool("tool_1", ToolConfig{
		Type:        ToolTypeClient,
		Name:        "show_order_status",
		Description: "Updated description.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ToolConfig.Description != "Updated description." {
		t.Fatalf("tool = %+v", got)
	}
}

func TestDeleteTool(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools/tool_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("force") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.DeleteTool("tool_1", WithForceDeleteTool()); err != nil {
		t.Fatal(err)
	}
}

func TestGetToolDependentAgents(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/tools/tool_1/dependent-agents" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_size") != "10" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"agents":[{"type":"available","access_level":"admin","created_at_unix_secs":1,"id":"agent_1","name":"Support"}],"has_more":false,"branches":[{"agent_id":"agent_1","agent_name":"Support","branch_id":"br_1","branch_name":"main","is_main":true}]}`))
	})

	got, err := client.GetToolDependentAgents("tool_1", WithToolDependentPageSize(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Agents) != 1 || got.Agents[0].ID != "agent_1" || got.Agents[0].Type != DependentAgentAvailable {
		t.Fatalf("agents = %+v", got.Agents)
	}
	if len(got.Branches) != 1 || !got.Branches[0].IsMain {
		t.Fatalf("branches = %+v", got.Branches)
	}
}

func TestToolError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if _, err := client.GetTool("missing"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := client.CreateTool(ToolConfig{Type: ToolTypeClient, Name: "x"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestToolsAll(t *testing.T) {
	tools := Tools{
		{ID: "a", ToolConfig: ToolConfig{Type: ToolTypeClient}},
		{ID: "b", ToolConfig: ToolConfig{Type: ToolTypeWebhook}},
	}

	var ids []string
	for tool := range tools.All() {
		ids = append(ids, tool.ID)
	}
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("ids = %v", ids)
	}

	var webhooks []string
	for tool := range tools.AllWith(func(tool *Tool) bool {
		return tool.ToolConfig.Type == ToolTypeWebhook
	}) {
		webhooks = append(webhooks, tool.ID)
	}
	if len(webhooks) != 1 || webhooks[0] != "b" {
		t.Fatalf("webhooks = %v", webhooks)
	}
}

func TestLiveTools(t *testing.T) {
	client := liveClient(t)
	created, err := client.CreateTool(ToolConfig{
		Type:        ToolTypeClient,
		Name:        "go_eleven_tool_test",
		Description: "Client tool created by the go-eleven SDK live test.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" {
		t.Fatal("empty tool id")
	}
	t.Cleanup(func() {
		_ = client.DeleteTool(created.ID, WithForceDeleteTool())
	})

	got, err := client.GetTool(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID || got.ToolConfig.Name != "go_eleven_tool_test" {
		t.Fatalf("get = %+v", got)
	}

	list, err := client.ListTools(WithToolSearch("go_eleven_tool_test"), WithToolPageSize(30))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for tool := range list.Tools.All() {
		if tool.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created tool missing from list: %+v", list)
	}

	updated, err := client.UpdateTool(created.ID, ToolConfig{
		Type:        ToolTypeClient,
		Name:        "go_eleven_tool_test",
		Description: "Updated by the go-eleven SDK live test.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.ToolConfig.Description != "Updated by the go-eleven SDK live test." {
		t.Fatalf("update = %+v", updated)
	}

	if err := client.DeleteTool(created.ID, WithForceDeleteTool()); err != nil {
		t.Fatal(err)
	}
}
