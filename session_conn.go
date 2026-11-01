package eleven

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const sessionReadLimit = 4 << 20

type conversationSession struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc

	events chan SessionEvent
	done   chan struct{}

	writeMu   sync.Mutex
	closeOnce sync.Once

	mu             sync.Mutex
	conversationID string
	err            error
}

func (c *client) StartConversation(agentID string, options ...StartConversationOption) (ConversationSession, error) {
	return c.StartConversationWithContext(context.Background(), agentID, options...)
}

func (c *client) StartConversationWithContext(ctx context.Context, agentID string, options ...StartConversationOption) (ConversationSession, error) {
	var opt startConversationOption
	for _, option := range options {
		option(&opt)
	}

	if opt.URL == "" && agentID == "" {
		return nil, fmt.Errorf("agentID is empty")
	}

	dialURL, sendKey, err := c.sessionDialURL(ctx, agentID, opt)
	if err != nil {
		return nil, err
	}

	header := make(http.Header)
	if sendKey {
		header.Set(xiAPIKeyHeader, c.apikey)
	}

	httpClient := *c.clnt
	httpClient.Timeout = 0

	sessionCtx, cancel := context.WithCancel(ctx)
	conn, resp, err := websocket.Dial(sessionCtx, dialURL, &websocket.DialOptions{
		HTTPClient: &httpClient,
		HTTPHeader: header,
	})
	if err != nil {
		cancel()
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("api response: %s: %s", resp.Status, body)
		}
		return nil, err
	}
	conn.SetReadLimit(sessionReadLimit)

	sess := &conversationSession{
		conn:   conn,
		ctx:    sessionCtx,
		cancel: cancel,
		events: make(chan SessionEvent, opt.eventBuffer()),
		done:   make(chan struct{}),
	}

	if err := sess.writeJSON(opt.initiation()); err != nil {
		cancel()
		_ = conn.Close(websocket.StatusGoingAway, "")
		return nil, err
	}

	go sess.readLoop()
	return sess, nil
}

func (c *client) sessionDialURL(ctx context.Context, agentID string, opt startConversationOption) (string, bool, error) {
	if opt.URL != "" {
		return opt.URL, false, nil
	}
	if opt.SignedURL {
		var signedOpts []SignedURLOption
		if opt.BranchID != "" {
			signedOpts = append(signedOpts, WithSignedURLBranchID(opt.BranchID))
		}
		if opt.VersionID != "" {
			signedOpts = append(signedOpts, WithSignedURLVersionID(opt.VersionID))
		}
		signed, err := c.GetSignedURLWithContext(ctx, agentID, signedOpts...)
		if err != nil {
			return "", false, err
		}
		if signed.SignedURL == "" {
			return "", false, fmt.Errorf("signed URL is empty")
		}
		return signed.SignedURL, false, nil
	}

	endpoint := websocketURL(c.base.JoinPath(conversationWSPath))
	endpoint.RawQuery = opt.query(agentID).Encode()
	return endpoint.String(), true, nil
}

func websocketURL(u *url.URL) *url.URL {
	cloned := *u
	switch cloned.Scheme {
	case "https":
		cloned.Scheme = "wss"
	case "http":
		cloned.Scheme = "ws"
	}
	return &cloned
}

func (s *conversationSession) Events() <-chan SessionEvent {
	return s.events
}

func (s *conversationSession) ConversationID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conversationID
}

func (s *conversationSession) SendUserMessage(text string) error {
	if text == "" {
		return fmt.Errorf("text is empty")
	}
	return s.writeJSON(map[string]any{
		"type": "user_message",
		"text": text,
	})
}

func (s *conversationSession) SendAudio(pcm []byte) error {
	if len(pcm) == 0 {
		return fmt.Errorf("audio is empty")
	}
	return s.writeJSON(map[string]any{
		"user_audio_chunk": base64.StdEncoding.EncodeToString(pcm),
	})
}

func (s *conversationSession) SendContextualUpdate(text string) error {
	if text == "" {
		return fmt.Errorf("text is empty")
	}
	return s.writeJSON(map[string]any{
		"type": "contextual_update",
		"text": text,
	})
}

func (s *conversationSession) SendUserActivity() error {
	return s.writeJSON(map[string]any{
		"type": "user_activity",
	})
}

func (s *conversationSession) SendClientToolResult(toolCallID string, result any, isError bool) error {
	if toolCallID == "" {
		return fmt.Errorf("toolCallID is empty")
	}
	return s.writeJSON(map[string]any{
		"type":         "client_tool_result",
		"tool_call_id": toolCallID,
		"result":       result,
		"is_error":     isError,
	})
}

func (s *conversationSession) SendFeedback(eventID int, feedback ConversationFeedback) error {
	if eventID == 0 {
		return fmt.Errorf("eventID is empty")
	}
	return s.writeJSON(map[string]any{
		"type":     "feedback",
		"event_id": eventID,
		"score":    feedback.String(),
	})
}

func (s *conversationSession) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		s.cancel()
		closeErr = s.conn.Close(websocket.StatusNormalClosure, "")
		<-s.done
	})
	return sessionCloseErr(closeErr)
}

func (s *conversationSession) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func (s *conversationSession) writeJSON(v any) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return wsjson.Write(s.ctx, s.conn, v)
}

func (s *conversationSession) readLoop() {
	defer close(s.done)
	defer close(s.events)
	defer s.conn.Close(websocket.StatusGoingAway, "")

	err := s.consume()
	s.mu.Lock()
	s.err = sessionCloseErr(err)
	s.mu.Unlock()
}

func (s *conversationSession) consume() error {
	for {
		_, body, err := s.conn.Read(s.ctx)
		if err != nil {
			return err
		}

		event, err := decodeSessionEvent(body)
		if err != nil {
			return err
		}

		if event.Type == SessionEventInitiation && event.Initiation != nil {
			s.mu.Lock()
			s.conversationID = event.Initiation.ConversationID
			s.mu.Unlock()
		}

		if event.Type == SessionEventPing && event.Ping != nil {
			if err := s.pong(*event.Ping); err != nil {
				return err
			}
		}

		select {
		case s.events <- event:
		case <-s.ctx.Done():
			return s.ctx.Err()
		}
	}
}

func (s *conversationSession) pong(event SessionPing) error {
	if event.PingMS > 0 {
		timer := time.NewTimer(time.Duration(event.PingMS) * time.Millisecond)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-s.ctx.Done():
			return s.ctx.Err()
		}
	}
	return s.writeJSON(map[string]any{
		"type":     "pong",
		"event_id": event.EventID,
	})
}

func sessionCloseErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	switch websocket.CloseStatus(err) {
	case websocket.StatusNormalClosure, websocket.StatusGoingAway, websocket.StatusNoStatusRcvd:
		return nil
	}
	return err
}

var _ ConversationSession = (*conversationSession)(nil)
