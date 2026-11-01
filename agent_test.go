package eleven

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAgentEnums(t *testing.T) {
	if AgentSortCallCount7d.String() != "call_count_7d" {
		t.Fatalf("sort = %s", AgentSortCallCount7d)
	}
	if AgentAccessRoleEditor.String() != "editor" {
		t.Fatalf("role = %s", AgentAccessRoleEditor)
	}
	if AgentAccessSourceWorkspaceAdmin.String() != "workspace_admin" {
		t.Fatalf("source = %s", AgentAccessSourceWorkspaceAdmin)
	}
}

func TestListAgentsQuery(t *testing.T) {
	var opt listAgentsOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithAgentPageSize(30)(&opt)
	WithAgentSearch("Support")(&opt)
	WithAgentArchived()(&opt)
	WithCreatedByUserID("@me")(&opt)
	WithAgentTags("support", "sales")(&opt)
	WithAgentSort(AgentSortName)(&opt)
	WithAgentSortDirection(SortAsc)(&opt)
	WithAgentCursor("cur_1")(&opt)

	q := opt.query()
	if q.Get("page_size") != "30" {
		t.Fatalf("page_size = %q", q.Get("page_size"))
	}
	if q.Get("search") != "Support" {
		t.Fatalf("search = %q", q.Get("search"))
	}
	if q.Get("archived") != "true" {
		t.Fatalf("archived = %q", q.Get("archived"))
	}
	if q.Get("created_by_user_id") != "@me" {
		t.Fatalf("created_by_user_id = %q", q.Get("created_by_user_id"))
	}
	if got := q["tags"]; len(got) != 2 || got[0] != "support" || got[1] != "sales" {
		t.Fatalf("tags = %v", got)
	}
	if q.Get("sort_by") != "name" {
		t.Fatalf("sort_by = %q", q.Get("sort_by"))
	}
	if q.Get("sort_direction") != "asc" {
		t.Fatalf("sort_direction = %q", q.Get("sort_direction"))
	}
	if q.Get("cursor") != "cur_1" {
		t.Fatalf("cursor = %q", q.Get("cursor"))
	}
}

func TestGetAgentQuery(t *testing.T) {
	var opt getAgentOption
	WithAgentVersionID("ver_1")(&opt)
	WithAgentBranchID("br_1")(&opt)
	WithoutProcedureGeneratedWorkflow()(&opt)

	q := opt.query()
	if q.Get("version_id") != "ver_1" {
		t.Fatalf("version_id = %q", q.Get("version_id"))
	}
	if q.Get("branch_id") != "br_1" {
		t.Fatalf("branch_id = %q", q.Get("branch_id"))
	}
	if q.Get("include_procedure_generated_workflow_elements") != "false" {
		t.Fatalf("include = %q", q.Get("include_procedure_generated_workflow_elements"))
	}
}

func TestCreateAgentRequestBody(t *testing.T) {
	var opt createAgentOption
	got, err := json.Marshal(opt.body(ConversationalConfig{
		Agent: &AgentConfig{FirstMessage: "Hi", Language: "en"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"conversation_config":{"agent":{"first_message":"Hi","language":"en"}}}` {
		t.Fatalf("minimal body = %s", got)
	}

	WithAgentName("Support Concierge")(&opt)
	WithCreateAgentTags("support")(&opt)
	WithPlatformSettings(json.RawMessage(`{"archived":false}`))(&opt)
	WithAgentWorkflow(json.RawMessage(`{"nodes":{}}`))(&opt)

	got, err = json.Marshal(opt.body(ConversationalConfig{
		TTS: &AgentTTSConfig{VoiceID: "voice-1", ModelID: "eleven_turbo_v2_5"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["name"] != "Support Concierge" {
		t.Fatalf("name = %v", payload["name"])
	}
	tags := payload["tags"].([]any)
	if len(tags) != 1 || tags[0] != "support" {
		t.Fatalf("tags = %v", tags)
	}
	settings := payload["platform_settings"].(map[string]any)
	if settings["archived"] != false {
		t.Fatalf("platform_settings = %v", settings)
	}
}

func TestUpdateAgentRequestBody(t *testing.T) {
	var opt updateAgentOption
	got, err := json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{}` {
		t.Fatalf("empty patch = %s", got)
	}

	WithUpdateAgentName("Support v2")(&opt)
	WithUpdateAgentTags("support", "vip")(&opt)
	WithVersionDescription("rename")(&opt)
	WithUpdateConversationConfig(ConversationalConfig{
		Agent: &AgentConfig{FirstMessage: "Hello"},
	})(&opt)

	got, err = json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["name"] != "Support v2" {
		t.Fatalf("name = %v", payload["name"])
	}
	if payload["version_description"] != "rename" {
		t.Fatalf("version_description = %v", payload["version_description"])
	}
	config := payload["conversation_config"].(map[string]any)
	agent := config["agent"].(map[string]any)
	if agent["first_message"] != "Hello" {
		t.Fatalf("conversation_config = %v", config)
	}
}

func TestAgentValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.GetAgent(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("get err = %v", err)
	}
	if _, err := client.UpdateAgent(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("update err = %v", err)
	}
	if err := client.DeleteAgent(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("delete err = %v", err)
	}
	if _, err := client.DuplicateAgent(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("duplicate err = %v", err)
	}
	if _, err := client.GetAgentLink(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("link err = %v", err)
	}
}

func TestListAgents(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("search") != "Support" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("page_size") != "20" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"agents":[{"agent_id":"agent_1","name":"Support","voice_id":"voice-1","tags":["support"],"created_at_unix_secs":1,"updated_at_unix_secs":2,"access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"archived":false}],"has_more":true,"next_cursor":"cur_2"}`))
	})

	got, err := client.ListAgents(WithAgentSearch("Support"), WithAgentPageSize(20))
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasMore || got.NextCursor != "cur_2" || len(got.Agents) != 1 {
		t.Fatalf("list = %+v", got)
	}
	if got.Agents[0].AgentID != "agent_1" || got.Agents[0].Name != "Support" {
		t.Fatalf("agent = %+v", got.Agents[0])
	}
}

func TestCreateAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/create" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		var payload createAgentRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Name != "Support Concierge" {
			t.Fatalf("name = %q", payload.Name)
		}
		if payload.ConversationConfig.TTS == nil || payload.ConversationConfig.TTS.VoiceID != "voice-1" {
			t.Fatalf("config = %+v", payload.ConversationConfig)
		}
		if payload.ConversationConfig.Agent == nil || payload.ConversationConfig.Agent.FirstMessage != "Hi" {
			t.Fatalf("agent = %+v", payload.ConversationConfig.Agent)
		}
		w.Write([]byte(`{"agent_id":"agent_1"}`))
	})

	got, err := client.CreateAgent(
		ConversationalConfig{
			Agent: &AgentConfig{FirstMessage: "Hi", Language: "en", Prompt: &AgentPrompt{Prompt: "Be concise.", LLM: "gemini-2.0-flash"}},
			TTS:   &AgentTTSConfig{VoiceID: "voice-1", ModelID: "eleven_turbo_v2_5"},
		},
		WithAgentName("Support Concierge"),
		WithCreateAgentTags("support"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentID != "agent_1" {
		t.Fatalf("result = %+v", got)
	}
}

func TestGetAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/agent_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("version_id") != "ver_1" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"agent_id":"agent_1","name":"Support","conversation_config":{"agent":{"first_message":"Hi","language":"en"},"tts":{"voice_id":"voice-1"}},"metadata":{"created_at_unix_secs":1,"updated_at_unix_secs":2},"tags":["support"],"version_id":"ver_1"}`))
	})

	got, err := client.GetAgent("agent_1", WithAgentVersionID("ver_1"))
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentID != "agent_1" || got.Name != "Support" || got.VersionID != "ver_1" {
		t.Fatalf("agent = %+v", got)
	}
	if got.ConversationConfig.Agent == nil || got.ConversationConfig.Agent.FirstMessage != "Hi" {
		t.Fatalf("config = %+v", got.ConversationConfig)
	}
}

func TestUpdateAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/agent_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("branch_id") != "br_1" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		var payload updateAgentRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Name == nil || *payload.Name != "Support v2" {
			t.Fatalf("payload = %+v", payload)
		}
		if payload.ConversationConfig != nil {
			t.Fatal("unexpected conversation_config")
		}
		w.Write([]byte(`{"agent_id":"agent_1","name":"Support v2","conversation_config":{},"metadata":{"created_at_unix_secs":1,"updated_at_unix_secs":3}}`))
	})

	got, err := client.UpdateAgent("agent_1", WithUpdateAgentName("Support v2"), WithUpdateAgentBranchID("br_1"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Support v2" {
		t.Fatalf("agent = %+v", got)
	}
}

func TestDeleteAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/agent_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("force") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.DeleteAgent("agent_1", WithForceDeleteAgent()); err != nil {
		t.Fatal(err)
	}
}

func TestDuplicateAgent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/agent_1/duplicate" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload duplicateAgentRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Name != "Support copy" {
			t.Fatalf("name = %q", payload.Name)
		}
		w.Write([]byte(`{"agent_id":"agent_2"}`))
	})

	got, err := client.DuplicateAgent("agent_1", WithDuplicateAgentName("Support copy"))
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentID != "agent_2" {
		t.Fatalf("result = %+v", got)
	}
}

func TestGetAgentLink(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/agents/agent_1/link" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{"agent_id":"agent_1","token":{"agent_id":"agent_1","conversation_token":"tok","purpose":"shareable_link","expiration_time_unix_secs":9}}`))
	})

	got, err := client.GetAgentLink("agent_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentID != "agent_1" || got.Token == nil || got.Token.ConversationToken != "tok" || got.Token.Purpose != "shareable_link" {
		t.Fatalf("link = %+v", got)
	}
}

func TestAgentError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if _, err := client.GetAgent("missing"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := client.CreateAgent(ConversationalConfig{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestAgentsAll(t *testing.T) {
	agents := Agents{
		{Name: "a", Archived: false},
		{Name: "b", Archived: true},
	}

	var names []string
	for agent := range agents.All() {
		names = append(names, agent.Name)
	}
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("names = %v", names)
	}

	var archived []string
	for agent := range agents.AllWith(func(a *AgentSummary) bool {
		return a.Archived
	}) {
		archived = append(archived, agent.Name)
	}
	if len(archived) != 1 || archived[0] != "b" {
		t.Fatalf("archived = %v", archived)
	}
}

func TestLiveAgents(t *testing.T) {
	client := liveClient(t)
	created, err := client.CreateAgent(
		ConversationalConfig{
			Agent: &AgentConfig{
				FirstMessage: "Hello from go eleven.",
				Language:     "en",
				Prompt:       &AgentPrompt{Prompt: "You are a concise test agent for the go-eleven SDK.", LLM: "gemini-2.0-flash"},
			},
			TTS: &AgentTTSConfig{
				VoiceID: "CwhRBWXzGAHq8TQ4Fs17",
				ModelID: "eleven_turbo_v2_5",
			},
		},
		WithAgentName("go-eleven-test"),
		WithCreateAgentTags("go-eleven"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if created.AgentID == "" {
		t.Fatal("empty agent id")
	}

	var extra []string
	t.Cleanup(func() {
		_ = client.DeleteAgent(created.AgentID, WithForceDeleteAgent())
		for _, id := range extra {
			_ = client.DeleteAgent(id, WithForceDeleteAgent())
		}
	})

	got, err := client.GetAgent(created.AgentID)
	if err != nil {
		t.Fatal(err)
	}
	if got.AgentID != created.AgentID || got.Name != "go-eleven-test" {
		t.Fatalf("get = %+v", got)
	}

	list, err := client.ListAgents(WithAgentSearch("go-eleven-test"), WithAgentPageSize(30))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for agent := range list.Agents.All() {
		if agent.AgentID == created.AgentID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("created agent missing from list")
	}

	updated, err := client.UpdateAgent(created.AgentID, WithUpdateAgentName("go-eleven-test-v2"))
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "go-eleven-test-v2" {
		t.Fatalf("updated name = %q", updated.Name)
	}

	link, err := client.GetAgentLink(created.AgentID)
	if err != nil {
		t.Fatal(err)
	}
	if link.AgentID != created.AgentID {
		t.Fatalf("link = %+v", link)
	}

	dup, err := client.DuplicateAgent(created.AgentID, WithDuplicateAgentName("go-eleven-test-copy"))
	if err != nil {
		t.Fatal(err)
	}
	if dup.AgentID == "" || dup.AgentID == created.AgentID {
		t.Fatalf("duplicate = %+v", dup)
	}
	extra = append(extra, dup.AgentID)
}
