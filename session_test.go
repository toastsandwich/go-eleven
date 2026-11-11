package eleven

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestSessionEnums(t *testing.T) {
	if SessionEventInitiation.String() != "conversation_initiation_metadata" {
		t.Fatalf("initiation = %s", SessionEventInitiation)
	}
	if SessionEventAgentChatResponsePart.String() != "agent_chat_response_part" {
		t.Fatalf("chat part = %s", SessionEventAgentChatResponsePart)
	}
	if SessionEventQueueStatus.String() != "queue_status" {
		t.Fatalf("queue = %s", SessionEventQueueStatus)
	}
	if SessionChatPartDelta.String() != "delta" {
		t.Fatalf("part = %s", SessionChatPartDelta)
	}
	if SessionQueueAdmitted.String() != "admitted" {
		t.Fatalf("queue status = %s", SessionQueueAdmitted)
	}
}

func TestSessionDialQuery(t *testing.T) {
	var opt startConversationOption
	q := opt.query("agent_1")
	if q.Get("agent_id") != "agent_1" {
		t.Fatalf("agent_id = %q", q.Get("agent_id"))
	}
	if q.Get("source") != conversationWSSource {
		t.Fatalf("source = %q", q.Get("source"))
	}

	WithSessionBranchID("br_1")(&opt)
	WithSessionVersionID("ver_1")(&opt)
	q = opt.query("agent_1")
	if q.Get("branch_id") != "br_1" || q.Get("version_id") != "ver_1" {
		t.Fatalf("query = %v", q)
	}
}

func TestSessionInitiationBody(t *testing.T) {
	var opt startConversationOption
	got, err := json.Marshal(opt.initiation())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"type":"conversation_initiation_client_data"}` {
		t.Fatalf("empty body = %s", got)
	}

	WithSessionUserID("user_1")(&opt)
	WithSessionConfigOverride(json.RawMessage(`{"agent":{"first_message":"hi"}}`))(&opt)
	WithSessionDynamicVariables(json.RawMessage(`{"name":"Ada"}`))(&opt)
	WithSessionExtraBody(json.RawMessage(`{"temperature":0.2}`))(&opt)
	got, err = json.Marshal(opt.initiation())
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["type"] != "conversation_initiation_client_data" || payload["user_id"] != "user_1" {
		t.Fatalf("payload = %s", got)
	}
}

func TestSessionEventBuffer(t *testing.T) {
	var opt startConversationOption
	if opt.eventBuffer() != defaultSessionEventBuffer {
		t.Fatalf("default = %d", opt.eventBuffer())
	}
	WithSessionEventBuffer(8)(&opt)
	if opt.eventBuffer() != 8 {
		t.Fatalf("buffer = %d", opt.eventBuffer())
	}
	WithSessionEventBuffer(0)(&opt)
	if opt.eventBuffer() != defaultSessionEventBuffer {
		t.Fatalf("zero = %d", opt.eventBuffer())
	}
}

func TestDecodeSessionEvent(t *testing.T) {
	pcm := []byte{0, 1, 2, 3}
	audioBody := fmt.Sprintf(
		`{"type":"audio","audio_event":{"audio_base_64":"%s","event_id":7}}`,
		base64.StdEncoding.EncodeToString(pcm),
	)
	event, err := decodeSessionEvent([]byte(audioBody))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != SessionEventAudio || event.Audio == nil || event.Audio.EventID != 7 {
		t.Fatalf("audio event = %+v", event)
	}
	if string(event.Audio.PCM) != string(pcm) {
		t.Fatalf("pcm = %v", event.Audio.PCM)
	}

	event, err = decodeSessionEvent([]byte(`{"type":"conversation_initiation_metadata","conversation_initiation_metadata_event":{"conversation_id":"conv_1","agent_output_audio_format":"pcm_16000","user_input_audio_format":"pcm_16000"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != SessionEventInitiation || event.Initiation.ConversationID != "conv_1" {
		t.Fatalf("initiation = %+v", event)
	}

	event, err = decodeSessionEvent([]byte(`{"type":"mcp_tool_call","mcp_tool_call":{"tool_name":"search"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != SessionEventUnknown {
		t.Fatalf("unknown type = %s", event.Type)
	}
	if string(event.Raw) == "" {
		t.Fatal("raw is empty")
	}

	event, err = decodeSessionEvent([]byte(`{"foo":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != SessionEventUnknown {
		t.Fatalf("empty type = %s", event.Type)
	}
}

func TestSessionValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.StartConversation(""); err == nil || err.Error() != "agentID is empty" {
		t.Fatalf("start err = %v", err)
	}

	s := &conversationSession{ctx: context.Background()}
	if err := s.SendUserMessage(""); err == nil || err.Error() != "text is empty" {
		t.Fatalf("message err = %v", err)
	}
	if err := s.SendContextualUpdate(""); err == nil || err.Error() != "text is empty" {
		t.Fatalf("update err = %v", err)
	}
	if err := s.SendAudio(nil); err == nil || err.Error() != "audio is empty" {
		t.Fatalf("audio err = %v", err)
	}
	if err := s.SendClientToolResult("", "ok", false); err == nil || err.Error() != "toolCallID is empty" {
		t.Fatalf("tool err = %v", err)
	}
	if err := s.SendFeedback(0, ConversationFeedbackLike); err == nil || err.Error() != "eventID is empty" {
		t.Fatalf("feedback err = %v", err)
	}
}

func TestStartConversation(t *testing.T) {
	gotInit := make(chan map[string]any, 1)
	gotUser := make(chan map[string]any, 1)
	gotPong := make(chan map[string]any, 1)

	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/conversation" {
			t.Errorf("path = %s", r.URL.Path)
			return
		}
		if r.URL.Query().Get("agent_id") != "agent_1" || r.URL.Query().Get("source") != conversationWSSource {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("branch_id") != "br_1" || r.URL.Query().Get("version_id") != "ver_1" {
			t.Errorf("query extras = %s", r.URL.RawQuery)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.CloseNow()

		var init map[string]any
		if err := wsjson.Read(r.Context(), conn, &init); err != nil {
			t.Error(err)
			return
		}
		gotInit <- init

		if err := wsjson.Write(r.Context(), conn, map[string]any{
			"type": "conversation_initiation_metadata",
			"conversation_initiation_metadata_event": map[string]any{
				"conversation_id":           "conv_live",
				"agent_output_audio_format": "pcm_16000",
				"user_input_audio_format":   "pcm_16000",
			},
		}); err != nil {
			t.Error(err)
			return
		}
		if err := wsjson.Write(r.Context(), conn, map[string]any{
			"type": "agent_response",
			"agent_response_event": map[string]any{
				"agent_response": "Hello from go eleven.",
				"event_id":       1,
			},
		}); err != nil {
			t.Error(err)
			return
		}

		var user map[string]any
		if err := wsjson.Read(r.Context(), conn, &user); err != nil {
			t.Error(err)
			return
		}
		gotUser <- user

		if err := wsjson.Write(r.Context(), conn, map[string]any{
			"type": "ping",
			"ping_event": map[string]any{
				"event_id": 9,
				"ping_ms":  5,
			},
		}); err != nil {
			t.Error(err)
			return
		}

		var pong map[string]any
		if err := wsjson.Read(r.Context(), conn, &pong); err != nil {
			t.Error(err)
			return
		}
		gotPong <- pong

		_, _, _ = conn.Read(r.Context())
	})

	session, err := client.StartConversation(
		"agent_1",
		WithSessionUserID("user_1"),
		WithSessionBranchID("br_1"),
		WithSessionVersionID("ver_1"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	sawInit, sawResponse, sawPing := false, false, false
	deadline := time.After(2 * time.Second)
	for !sawInit || !sawResponse || !sawPing {
		select {
		case event, ok := <-session.Events():
			if !ok {
				t.Fatal("events closed")
			}
			switch event.Type {
			case SessionEventInitiation:
				if event.Initiation.ConversationID != "conv_live" {
					t.Fatalf("id = %s", event.Initiation.ConversationID)
				}
				if session.ConversationID() != "conv_live" {
					t.Fatalf("session id = %s", session.ConversationID())
				}
				sawInit = true
			case SessionEventAgentResponse:
				if event.Response.Text != "Hello from go eleven." {
					t.Fatalf("response = %s", event.Response.Text)
				}
				sawResponse = true
				if err := session.SendUserMessage("What are your hours?"); err != nil {
					t.Fatal(err)
				}
			case SessionEventPing:
				if event.Ping.EventID != 9 {
					t.Fatalf("ping = %+v", event.Ping)
				}
				sawPing = true
			default:
				t.Fatalf("unexpected event %s", event.Type)
			}
		case <-deadline:
			t.Fatalf("timeout init=%t response=%t ping=%t", sawInit, sawResponse, sawPing)
		}
	}

	select {
	case init := <-gotInit:
		if init["type"] != "conversation_initiation_client_data" || init["user_id"] != "user_1" {
			t.Fatalf("init = %v", init)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for initiation frame")
	}
	select {
	case user := <-gotUser:
		if user["type"] != "user_message" || user["text"] != "What are your hours?" {
			t.Fatalf("user = %v", user)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for user message")
	}
	select {
	case pong := <-gotPong:
		if pong["type"] != "pong" || pong["event_id"] != float64(9) {
			t.Fatalf("pong = %v", pong)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for pong")
	}

	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if err := session.Err(); err != nil {
		t.Fatalf("session err = %v", err)
	}
}

func TestStartConversationSendFrames(t *testing.T) {
	got := make(chan map[string]any, 8)
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.CloseNow()

		var init map[string]any
		if err := wsjson.Read(r.Context(), conn, &init); err != nil {
			t.Error(err)
			return
		}
		got <- init

		if err := wsjson.Write(r.Context(), conn, map[string]any{
			"type": "conversation_initiation_metadata",
			"conversation_initiation_metadata_event": map[string]any{
				"conversation_id": "conv_frames",
			},
		}); err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < 5; i++ {
			var frame map[string]any
			if err := wsjson.Read(r.Context(), conn, &frame); err != nil {
				t.Error(err)
				return
			}
			got <- frame
		}
		_, _, _ = conn.Read(r.Context())
	})

	session, err := client.StartConversation("agent_1")
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	waitEvent(t, session, SessionEventInitiation)

	if err := session.SendAudio([]byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	if err := session.SendContextualUpdate("user is looking at checkout"); err != nil {
		t.Fatal(err)
	}
	if err := session.SendUserActivity(); err != nil {
		t.Fatal(err)
	}
	if err := session.SendClientToolResult("tool_1", "ok", false); err != nil {
		t.Fatal(err)
	}
	if err := session.SendFeedback(4, ConversationFeedbackLike); err != nil {
		t.Fatal(err)
	}

	frames := make([]map[string]any, 0, 6)
	for i := 0; i < 6; i++ {
		select {
		case frame := <-got:
			frames = append(frames, frame)
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout after %d frames", len(frames))
		}
	}

	if frames[0]["type"] != "conversation_initiation_client_data" {
		t.Fatalf("init = %v", frames[0])
	}
	if frames[1]["user_audio_chunk"] != base64.StdEncoding.EncodeToString([]byte{1, 2, 3}) {
		t.Fatalf("audio = %v", frames[1])
	}
	if frames[2]["type"] != "contextual_update" || frames[2]["text"] != "user is looking at checkout" {
		t.Fatalf("update = %v", frames[2])
	}
	if frames[3]["type"] != "user_activity" {
		t.Fatalf("activity = %v", frames[3])
	}
	if frames[4]["type"] != "client_tool_result" || frames[4]["tool_call_id"] != "tool_1" || frames[4]["is_error"] != false {
		t.Fatalf("tool = %v", frames[4])
	}
	if frames[5]["type"] != "feedback" || frames[5]["event_id"] != float64(4) || frames[5]["score"] != "like" {
		t.Fatalf("feedback = %v", frames[5])
	}
}

func TestStartConversationWithSessionURL(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/conversation" {
			t.Errorf("path = %s", r.URL.Path)
			return
		}
		if got := r.Header.Get(xiAPIKeyHeader); got != "" {
			t.Errorf("unexpected api key %q", got)
		}
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.CloseNow()
		var init map[string]any
		if err := wsjson.Read(r.Context(), conn, &init); err != nil {
			t.Error(err)
			return
		}
		_ = wsjson.Write(r.Context(), conn, map[string]any{
			"type": "conversation_initiation_metadata",
			"conversation_initiation_metadata_event": map[string]any{
				"conversation_id": "conv_url",
			},
		})
		_, _, _ = conn.Read(r.Context())
	})

	wsURL := websocketURL(client.base.JoinPath(conversationWSPath)).String()
	session, err := client.StartConversation("", WithSessionURL(wsURL))
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	event := waitEvent(t, session, SessionEventInitiation)
	if event.Initiation.ConversationID != "conv_url" {
		t.Fatalf("id = %s", event.Initiation.ConversationID)
	}
}

func TestStartConversationSignedURL(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/convai/conversation/get-signed-url":
			if r.URL.Query().Get("agent_id") != "agent_1" {
				t.Errorf("signed query = %s", r.URL.RawQuery)
			}
			requireHeader(t, r, xiAPIKeyHeader, "test-key")
			fmt.Fprintf(w, `{"signed_url":"ws://%s/v1/convai/conversation?agent_id=agent_1&conversation_signature=sig"}`, r.Host)
		case "/v1/convai/conversation":
			if got := r.Header.Get(xiAPIKeyHeader); got != "" {
				t.Errorf("unexpected api key %q", got)
			}
			conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.CloseNow()
			var init map[string]any
			if err := wsjson.Read(r.Context(), conn, &init); err != nil {
				t.Error(err)
				return
			}
			_ = wsjson.Write(r.Context(), conn, map[string]any{
				"type": "conversation_initiation_metadata",
				"conversation_initiation_metadata_event": map[string]any{
					"conversation_id": "conv_signed",
				},
			})
			_, _, _ = conn.Read(r.Context())
		default:
			t.Errorf("path = %s", r.URL.Path)
		}
	})

	session, err := client.StartConversation("agent_1", WithSessionSignedURL())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	event := waitEvent(t, session, SessionEventInitiation)
	if event.Initiation.ConversationID != "conv_signed" {
		t.Fatalf("id = %s", event.Initiation.ConversationID)
	}
}

func TestLiveConversationSession(t *testing.T) {
	client := liveClient(t)
	created, err := client.CreateAgent(
		ConversationalConfig{
			Agent: &AgentConfig{
				FirstMessage: "Hello from go eleven.",
				Language:     "en",
				Prompt:       &AgentPrompt{Prompt: "You are a concise test agent for the go-eleven SDK. Reply in one short sentence.", LLM: "gemini-2.0-flash"},
			},
			TTS: &AgentTTSConfig{
				VoiceID: "CwhRBWXzGAHq8TQ4Fs17",
				ModelID: "eleven_turbo_v2",
			},
		},
		WithAgentName("go-eleven-session-test"),
		WithCreateAgentTags("go-eleven"),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteAgent(created.AgentID, WithForceDeleteAgent())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	session, err := client.StartConversationWithContext(ctx, created.AgentID, WithSessionSignedURL())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if err := session.SendUserMessage("Say hello in one word."); err != nil {
		t.Fatal(err)
	}

	sawSpeech := false
	for !sawSpeech {
		select {
		case event, ok := <-session.Events():
			if !ok {
				t.Fatalf("events closed: %v", session.Err())
			}
			switch event.Type {
			case SessionEventAgentResponse, SessionEventAudio:
				sawSpeech = true
			case SessionEventClientError:
				t.Fatalf("client error: %+v", event.Error)
			}
		case <-ctx.Done():
			t.Fatalf("timeout waiting for agent speech: %v", session.Err())
		}
	}

	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
}

func waitEvent(t *testing.T, session ConversationSession, want SessionEventType) SessionEvent {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case event, ok := <-session.Events():
			if !ok {
				t.Fatalf("events closed waiting for %s: %v", want, session.Err())
			}
			if event.Type == want {
				return event
			}
		case <-deadline:
			t.Fatalf("timeout waiting for %s", want)
		}
	}
}
