package eleven

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestConversationEnums(t *testing.T) {
	if ConversationStatusInProgress.String() != "in-progress" {
		t.Fatalf("status = %s", ConversationStatusInProgress)
	}
	if CallSuccessfulSuccess.String() != "success" {
		t.Fatalf("call_successful = %s", CallSuccessfulSuccess)
	}
	if ConversationDirectionOutbound.String() != "outbound" {
		t.Fatalf("direction = %s", ConversationDirectionOutbound)
	}
	if ConversationFeedbackDislike.String() != "dislike" {
		t.Fatalf("feedback = %s", ConversationFeedbackDislike)
	}
	if ConversationSummaryInclude.String() != "include" {
		t.Fatalf("summary_mode = %s", ConversationSummaryInclude)
	}
}

func TestListConversationsQuery(t *testing.T) {
	var opt listConversationsOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithConversationCursor("cur_1")(&opt)
	WithConversationAgentID("agent_1")(&opt)
	WithConversationPageSize(30)(&opt)
	WithConversationSortDirection(SortAsc)(&opt)
	WithCallSuccessful(CallSuccessfulSuccess)(&opt)
	WithCallStartAfterUnix(10)(&opt)
	WithCallStartBeforeUnix(20)(&opt)
	WithConversationUserID("user_1")(&opt)
	WithConversationBranchID("br_1")(&opt)
	WithConversationVersionID("ver_1")(&opt)
	WithSummaryMode(ConversationSummaryInclude)(&opt)

	q := opt.query()
	if q.Get("cursor") != "cur_1" {
		t.Fatalf("cursor = %q", q.Get("cursor"))
	}
	if q.Get("agent_id") != "agent_1" {
		t.Fatalf("agent_id = %q", q.Get("agent_id"))
	}
	if q.Get("page_size") != "30" {
		t.Fatalf("page_size = %q", q.Get("page_size"))
	}
	if q.Get("sort_direction") != "asc" {
		t.Fatalf("sort_direction = %q", q.Get("sort_direction"))
	}
	if q.Get("call_successful") != "success" {
		t.Fatalf("call_successful = %q", q.Get("call_successful"))
	}
	if q.Get("call_start_after_unix") != "10" {
		t.Fatalf("call_start_after_unix = %q", q.Get("call_start_after_unix"))
	}
	if q.Get("call_start_before_unix") != "20" {
		t.Fatalf("call_start_before_unix = %q", q.Get("call_start_before_unix"))
	}
	if q.Get("user_id") != "user_1" {
		t.Fatalf("user_id = %q", q.Get("user_id"))
	}
	if q.Get("branch_id") != "br_1" {
		t.Fatalf("branch_id = %q", q.Get("branch_id"))
	}
	if q.Get("version_id") != "ver_1" {
		t.Fatalf("version_id = %q", q.Get("version_id"))
	}
	if q.Get("summary_mode") != "include" {
		t.Fatalf("summary_mode = %q", q.Get("summary_mode"))
	}
}

func TestSignedURLQuery(t *testing.T) {
	var opt signedURLOption
	q := opt.query("agent_1")
	if q.Get("agent_id") != "agent_1" {
		t.Fatalf("agent_id = %q", q.Get("agent_id"))
	}

	WithSignedURLConversationID()(&opt)
	WithSignedURLBranchID("br_1")(&opt)
	WithSignedURLVersionID("ver_1")(&opt)
	WithSignedURLEnvironment("staging")(&opt)
	q = opt.query("agent_1")
	if q.Get("include_conversation_id") != "true" {
		t.Fatalf("include_conversation_id = %q", q.Get("include_conversation_id"))
	}
	if q.Get("branch_id") != "br_1" || q.Get("version_id") != "ver_1" || q.Get("environment") != "staging" {
		t.Fatalf("query = %v", q)
	}
}

func TestConversationTokenQuery(t *testing.T) {
	var opt conversationTokenOption
	WithParticipantName("web-user")(&opt)
	WithTokenBranchID("br_1")(&opt)
	WithTokenVersionID("ver_1")(&opt)
	WithTokenEnvironment("staging")(&opt)
	q := opt.query("agent_1")
	if q.Get("agent_id") != "agent_1" || q.Get("participant_name") != "web-user" {
		t.Fatalf("query = %v", q)
	}
}

func TestRealtimeClientSecretBody(t *testing.T) {
	var opt clientSecretOption
	got, err := json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{}` {
		t.Fatalf("empty body = %s", got)
	}

	WithRealtimeExpiresAfter(600)(&opt)
	WithRealtimeSession(json.RawMessage(`{"type":"realtime"}`))(&opt)
	got, err = json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	expires := payload["expires_after"].(map[string]any)
	if expires["seconds"] != float64(600) || expires["anchor"] != "created_at" {
		t.Fatalf("expires_after = %v", expires)
	}
}

func TestConversationValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.GetConversation(""); err == nil || err.Error() != "conversationID is empty" {
		t.Fatalf("get err = %v", err)
	}
	if _, err := client.GetConversationSummary(""); err == nil || err.Error() != "conversationID is empty" {
		t.Fatalf("summary err = %v", err)
	}
	if err := client.DeleteConversation(""); err == nil || err.Error() != "conversationID is empty" {
		t.Fatalf("delete err = %v", err)
	}
	if _, err := client.GetConversationAudio(""); err == nil || err.Error() != "conversationID is empty" {
		t.Fatalf("audio err = %v", err)
	}
	if _, err := client.GetSignedURL(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("signed url err = %v", err)
	}
	if _, err := client.GetConversationToken(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("token err = %v", err)
	}
	if err := client.SendConversationFeedback("", ConversationFeedbackLike); err == nil || err.Error() != "conversationID is empty" {
		t.Fatalf("feedback id err = %v", err)
	}
	if err := client.SendConversationFeedback("conv_1", ""); err == nil || err.Error() != "feedback is empty" {
		t.Fatalf("feedback empty err = %v", err)
	}
}

func TestListConversations(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/conversations" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("agent_id") != "agent_1" || r.URL.Query().Get("page_size") != "30" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"conversations":[{"agent_id":"agent_1","conversation_id":"conv_1","start_time_unix_secs":1,"call_duration_secs":12,"message_count":4,"status":"done","call_successful":"success","agent_name":"Support"}],"has_more":true,"next_cursor":"cur_2"}`))
	})

	got, err := client.ListConversations(WithConversationAgentID("agent_1"), WithConversationPageSize(30))
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasMore || got.NextCursor != "cur_2" || len(got.Conversations) != 1 {
		t.Fatalf("list = %+v", got)
	}
	if got.Conversations[0].ConversationID != "conv_1" || got.Conversations[0].Status != ConversationStatusDone {
		t.Fatalf("item = %+v", got.Conversations[0])
	}
}

func TestGetConversation(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/conversations/conv_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{"agent_id":"agent_1","conversation_id":"conv_1","status":"done","has_audio":true,"has_user_audio":true,"has_response_audio":true,"has_auxiliary_audio":false,"metadata":{"start_time_unix_secs":1,"call_duration_secs":12,"termination_reason":"end_call"},"transcript":[{"role":"agent","message":"Hi","time_in_call_secs":0},{"role":"user","message":"Hello","time_in_call_secs":2}],"analysis":{"call_successful":"success","transcript_summary":"A greeting."}}`))
	})

	got, err := client.GetConversation("conv_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ConversationID != "conv_1" || !got.HasAudio || got.Metadata.CallDurationSecs != 12 {
		t.Fatalf("conversation = %+v", got)
	}
	if len(got.Transcript) != 2 || got.Transcript[1].Message != "Hello" {
		t.Fatalf("transcript = %+v", got.Transcript)
	}
	if got.Analysis == nil || got.Analysis.CallSuccessful != CallSuccessfulSuccess {
		t.Fatalf("analysis = %+v", got.Analysis)
	}
}

func TestGetConversationSummary(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/conversations/conv_1/summary" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("max_messages") != "20" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"conversation_id":"conv_1","agent_id":"agent_1","status":"done","message_count":2,"note":"fetch full transcript","call_summary_title":"Greeting","transcript_summary":"They said hello.","call_successful":"success","messages":[{"role":"agent","message":"Hi"}],"messages_omitted":false}`))
	})

	got, err := client.GetConversationSummary("conv_1", WithMaxMessages(20))
	if err != nil {
		t.Fatal(err)
	}
	if got.CallSummaryTitle != "Greeting" || len(got.Messages) != 1 {
		t.Fatalf("summary = %+v", got)
	}
}

func TestDeleteConversation(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/conversations/conv_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.DeleteConversation("conv_1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetConversationAudio(t *testing.T) {
	audio := []byte("ID3fake-mp3")
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/conversations/conv_1/audio" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write(audio)
	})

	got, err := client.GetConversationAudio("conv_1")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(audio) {
		t.Fatalf("audio = %q", got)
	}
}

func TestGetSignedURL(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/conversation/get-signed-url" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("agent_id") != "agent_1" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("include_conversation_id") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"signed_url":"wss://api.elevenlabs.io/v1/convai/conversation?token=abc"}`))
	})

	got, err := client.GetSignedURL("agent_1", WithSignedURLConversationID())
	if err != nil {
		t.Fatal(err)
	}
	if got.SignedURL == "" {
		t.Fatalf("signed url = %+v", got)
	}
}

func TestGetConversationToken(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/conversation/token" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("agent_id") != "agent_1" || r.URL.Query().Get("participant_name") != "web-user" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"token":"tok","conversation_id":"conv_1"}`))
	})

	got, err := client.GetConversationToken("agent_1", WithParticipantName("web-user"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Token != "tok" || got.ConversationID != "conv_1" {
		t.Fatalf("token = %+v", got)
	}
}

func TestSendConversationFeedback(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/conversations/conv_1/feedback" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload conversationFeedbackRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Feedback != "like" {
			t.Fatalf("feedback = %q", payload.Feedback)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.SendConversationFeedback("conv_1", ConversationFeedbackLike); err != nil {
		t.Fatal(err)
	}
}

func TestCreateRealtimeClientSecret(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/realtime/client_secrets" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload createClientSecretRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.ExpiresAfter == nil || payload.ExpiresAfter.Seconds != 600 {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"value":"secret","expires_at":1}`))
	})

	got, err := client.CreateRealtimeClientSecret(WithRealtimeExpiresAfter(600))
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != "secret" {
		t.Fatalf("secret = %+v", got)
	}
}

func TestConversationError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if _, err := client.GetConversation("missing"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := client.GetSignedURL("missing"); err == nil {
		t.Fatal("expected error")
	}
}

func TestConversationListAll(t *testing.T) {
	items := ConversationListItems{
		{ConversationID: "a", Status: ConversationStatusDone},
		{ConversationID: "b", Status: ConversationStatusFailed},
	}
	var ids []string
	for conv := range items.All() {
		ids = append(ids, conv.ConversationID)
	}
	if len(ids) != 2 {
		t.Fatalf("ids = %v", ids)
	}
	var failed []string
	for conv := range items.AllWith(func(c *ConversationListItem) bool {
		return c.Status == ConversationStatusFailed
	}) {
		failed = append(failed, conv.ConversationID)
	}
	if len(failed) != 1 || failed[0] != "b" {
		t.Fatalf("failed = %v", failed)
	}
}

func TestLiveConversations(t *testing.T) {
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
		WithAgentName("go-eleven-conversation-test"),
		WithCreateAgentTags("go-eleven"),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteAgent(created.AgentID, WithForceDeleteAgent())
	})

	signed, err := client.GetSignedURL(created.AgentID)
	if err != nil {
		t.Fatal(err)
	}
	if signed.SignedURL == "" {
		t.Fatal("empty signed url")
	}

	token, err := client.GetConversationToken(created.AgentID, WithParticipantName("go-eleven"))
	if err != nil {
		t.Fatal(err)
	}
	if token.Token == "" {
		t.Fatal("empty conversation token")
	}

	list, err := client.ListConversations(WithConversationAgentID(created.AgentID), WithConversationPageSize(30))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("conversations=%d signed=%t token=%t", len(list.Conversations), signed.SignedURL != "", token.Token != "")
}
