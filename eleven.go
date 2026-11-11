package eleven

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Client interface {
	ListModels() (Models, error)
	ListModelsWithContext(context.Context) (Models, error)
	SearchVoices(options ...SearchVoicesOption) (VoiceSearchResult, error)
	SearchVoicesWithContext(ctx context.Context, options ...SearchVoicesOption) (VoiceSearchResult, error)
	TextToSpeech(text string, voiceId string, options ...TextToSpeechOption) ([]byte, error)
	TextToSpeechWithContext(ctx context.Context, text string, voiceId string, options ...TextToSpeechOption) ([]byte, error)
	TextToSpeechStream(text string, voiceId string, options ...TextToSpeechOption) (io.ReadCloser, error)
	TextToSpeechStreamWithContext(ctx context.Context, text string, voiceId string, options ...TextToSpeechOption) (io.ReadCloser, error)
	SpeechToSpeech(audio io.Reader, voiceId string, options ...SpeechToSpeechOption) ([]byte, error)
	SpeechToSpeechWithContext(ctx context.Context, audio io.Reader, voiceId string, options ...SpeechToSpeechOption) ([]byte, error)
	SpeechToSpeechStream(audio io.Reader, voiceId string, options ...SpeechToSpeechOption) (io.ReadCloser, error)
	SpeechToSpeechStreamWithContext(ctx context.Context, audio io.Reader, voiceId string, options ...SpeechToSpeechOption) (io.ReadCloser, error)
	SpeechToText(audio io.Reader, options ...SpeechToTextOption) (SpeechToTextResult, error)
	SpeechToTextWithContext(ctx context.Context, audio io.Reader, options ...SpeechToTextOption) (SpeechToTextResult, error)
	GetTranscript(transcriptionID string) (SpeechToTextResult, error)
	GetTranscriptWithContext(ctx context.Context, transcriptionID string) (SpeechToTextResult, error)
	DeleteTranscript(transcriptionID string) error
	DeleteTranscriptWithContext(ctx context.Context, transcriptionID string) error
	ListAgents(options ...ListAgentsOption) (AgentList, error)
	ListAgentsWithContext(ctx context.Context, options ...ListAgentsOption) (AgentList, error)
	CreateAgent(config ConversationalConfig, options ...CreateAgentOption) (CreateAgentResult, error)
	CreateAgentWithContext(ctx context.Context, config ConversationalConfig, options ...CreateAgentOption) (CreateAgentResult, error)
	CreateOrUpdateAgent(config ConversationalConfig, options ...CreateAgentOption) (CreateOrUpdateAgentResult, error)
	CreateOrUpdateAgentWithContext(ctx context.Context, config ConversationalConfig, options ...CreateAgentOption) (CreateOrUpdateAgentResult, error)
	GetAgent(agentID string, options ...GetAgentOption) (Agent, error)
	GetAgentWithContext(ctx context.Context, agentID string, options ...GetAgentOption) (Agent, error)
	UpdateAgent(agentID string, options ...UpdateAgentOption) (Agent, error)
	UpdateAgentWithContext(ctx context.Context, agentID string, options ...UpdateAgentOption) (Agent, error)
	DeleteAgent(agentID string, options ...DeleteAgentOption) error
	DeleteAgentWithContext(ctx context.Context, agentID string, options ...DeleteAgentOption) error
	DuplicateAgent(agentID string, options ...DuplicateAgentOption) (CreateAgentResult, error)
	DuplicateAgentWithContext(ctx context.Context, agentID string, options ...DuplicateAgentOption) (CreateAgentResult, error)
	GetAgentLink(agentID string) (AgentLink, error)
	GetAgentLinkWithContext(ctx context.Context, agentID string) (AgentLink, error)
	ListConversations(options ...ListConversationsOption) (ConversationList, error)
	ListConversationsWithContext(ctx context.Context, options ...ListConversationsOption) (ConversationList, error)
	GetConversation(conversationID string) (Conversation, error)
	GetConversationWithContext(ctx context.Context, conversationID string) (Conversation, error)
	GetConversationSummary(conversationID string, options ...GetConversationSummaryOption) (ConversationSummary, error)
	GetConversationSummaryWithContext(ctx context.Context, conversationID string, options ...GetConversationSummaryOption) (ConversationSummary, error)
	DeleteConversation(conversationID string) error
	DeleteConversationWithContext(ctx context.Context, conversationID string) error
	GetConversationAudio(conversationID string) ([]byte, error)
	GetConversationAudioWithContext(ctx context.Context, conversationID string) ([]byte, error)
	GetSignedURL(agentID string, options ...SignedURLOption) (SignedURL, error)
	GetSignedURLWithContext(ctx context.Context, agentID string, options ...SignedURLOption) (SignedURL, error)
	GetConversationToken(agentID string, options ...ConversationTokenOption) (ConversationToken, error)
	GetConversationTokenWithContext(ctx context.Context, agentID string, options ...ConversationTokenOption) (ConversationToken, error)
	SendConversationFeedback(conversationID string, feedback ConversationFeedback) error
	SendConversationFeedbackWithContext(ctx context.Context, conversationID string, feedback ConversationFeedback) error
	CreateRealtimeClientSecret(options ...ClientSecretOption) (RealtimeClientSecret, error)
	CreateRealtimeClientSecretWithContext(ctx context.Context, options ...ClientSecretOption) (RealtimeClientSecret, error)
	StartConversation(agentID string, options ...StartConversationOption) (ConversationSession, error)
	StartConversationWithContext(ctx context.Context, agentID string, options ...StartConversationOption) (ConversationSession, error)
	ListTools(options ...ListToolsOption) (ToolList, error)
	ListToolsWithContext(ctx context.Context, options ...ListToolsOption) (ToolList, error)
	CreateTool(config ToolConfig, options ...CreateToolOption) (Tool, error)
	CreateToolWithContext(ctx context.Context, config ToolConfig, options ...CreateToolOption) (Tool, error)
	GetTool(toolID string, options ...GetToolOption) (Tool, error)
	GetToolWithContext(ctx context.Context, toolID string, options ...GetToolOption) (Tool, error)
	UpdateTool(toolID string, config ToolConfig, options ...UpdateToolOption) (Tool, error)
	UpdateToolWithContext(ctx context.Context, toolID string, config ToolConfig, options ...UpdateToolOption) (Tool, error)
	DeleteTool(toolID string, options ...DeleteToolOption) error
	DeleteToolWithContext(ctx context.Context, toolID string, options ...DeleteToolOption) error
	GetToolDependentAgents(toolID string, options ...ToolDependentOption) (DependentAgentList, error)
	GetToolDependentAgentsWithContext(ctx context.Context, toolID string, options ...ToolDependentOption) (DependentAgentList, error)
	ListKnowledgeBase(options ...ListKnowledgeOption) (KnowledgeList, error)
	ListKnowledgeBaseWithContext(ctx context.Context, options ...ListKnowledgeOption) (KnowledgeList, error)
	CreateKnowledgeBaseText(text string, options ...CreateKnowledgeTextOption) (CreateKnowledgeResult, error)
	CreateKnowledgeBaseTextWithContext(ctx context.Context, text string, options ...CreateKnowledgeTextOption) (CreateKnowledgeResult, error)
	CreateKnowledgeBaseURL(pageURL string, options ...CreateKnowledgeURLOption) (CreateKnowledgeResult, error)
	CreateKnowledgeBaseURLWithContext(ctx context.Context, pageURL string, options ...CreateKnowledgeURLOption) (CreateKnowledgeResult, error)
	CreateKnowledgeBaseFile(file io.Reader, options ...CreateKnowledgeFileOption) (CreateKnowledgeResult, error)
	CreateKnowledgeBaseFileWithContext(ctx context.Context, file io.Reader, options ...CreateKnowledgeFileOption) (CreateKnowledgeResult, error)
	GetKnowledgeBase(documentationID string, options ...GetKnowledgeOption) (KnowledgeDocument, error)
	GetKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...GetKnowledgeOption) (KnowledgeDocument, error)
	UpdateKnowledgeBase(documentationID string, options ...UpdateKnowledgeOption) (KnowledgeDocument, error)
	UpdateKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...UpdateKnowledgeOption) (KnowledgeDocument, error)
	DeleteKnowledgeBase(documentationID string, options ...DeleteKnowledgeOption) error
	DeleteKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...DeleteKnowledgeOption) error
	GetKnowledgeBaseContent(documentationID string) (KnowledgeContent, error)
	GetKnowledgeBaseContentWithContext(ctx context.Context, documentationID string) (KnowledgeContent, error)
	GetKnowledgeBaseDependentAgents(documentationID string, options ...KnowledgeDependentOption) (DependentAgentList, error)
	GetKnowledgeBaseDependentAgentsWithContext(ctx context.Context, documentationID string, options ...KnowledgeDependentOption) (DependentAgentList, error)
	ComputeRAGIndex(documentationID string, model EmbeddingModel) (RAGIndex, error)
	ComputeRAGIndexWithContext(ctx context.Context, documentationID string, model EmbeddingModel) (RAGIndex, error)
	GetKnowledgeBaseChunk(documentationID string, chunkID string, options ...GetKnowledgeChunkOption) (KnowledgeChunk, error)
	GetKnowledgeBaseChunkWithContext(ctx context.Context, documentationID string, chunkID string, options ...GetKnowledgeChunkOption) (KnowledgeChunk, error)
	ListWebhooks(options ...ListWebhooksOption) (WebhookList, error)
	ListWebhooksWithContext(ctx context.Context, options ...ListWebhooksOption) (WebhookList, error)
	CreateWebhook(name string, webhookURL string, options ...CreateWebhookOption) (CreateWebhookResult, error)
	CreateWebhookWithContext(ctx context.Context, name string, webhookURL string, options ...CreateWebhookOption) (CreateWebhookResult, error)
	UpdateWebhook(webhookID string, name string, isDisabled bool, options ...UpdateWebhookOption) (WebhookStatusResult, error)
	UpdateWebhookWithContext(ctx context.Context, webhookID string, name string, isDisabled bool, options ...UpdateWebhookOption) (WebhookStatusResult, error)
	DeleteWebhook(webhookID string) (WebhookStatusResult, error)
	DeleteWebhookWithContext(ctx context.Context, webhookID string) (WebhookStatusResult, error)
}

func createHTTPClient() *httpClient {
	url, err := url.Parse(parentEndpoint)
	if err != nil {
		panic(err)
	}

	clnt := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &httpClient{
		base: url,
		clnt: clnt,
	}
}

type httpClient struct {
	base *url.URL
	clnt *http.Client
}

type client struct {
	*httpClient
	apikey string
}

func NewClient(apiKey string) Client {
	if apiKey == "" {
		panic("apikey is empty")
	}
	return &client{
		httpClient: createHTTPClient(),
		apikey:     apiKey,
	}
}

func (c *client) newJSONRequest(ctx context.Context, method string, endpoint *url.URL, payload any) (*http.Request, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: method,
		URL:    endpoint,
		Header: make(http.Header),
		Body:   io.NopCloser(bytes.NewReader(b)),
	}
	req.ContentLength = int64(len(b))
	req = req.WithContext(ctx)
	req.Header.Add(xiAPIKeyHeader, c.apikey)
	req.Header.Add("Content-Type", "Application/json")
	return req, nil
}

func (c *client) doPOST(ctx context.Context, endpoint *url.URL, payload any) ([]byte, error) {
	return c.doJSON(ctx, http.MethodPost, endpoint, payload)
}

func (c *client) doPATCH(ctx context.Context, endpoint *url.URL, payload any) ([]byte, error) {
	return c.doJSON(ctx, http.MethodPatch, endpoint, payload)
}

func (c *client) doJSON(ctx context.Context, method string, endpoint *url.URL, payload any) ([]byte, error) {
	req, err := c.newJSONRequest(ctx, method, endpoint, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.clnt.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api response: %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *client) doPOSTStream(ctx context.Context, endpoint *url.URL, payload any) (io.ReadCloser, error) {
	req, err := c.newJSONRequest(ctx, http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}
	return c.doStream(req)
}

func (c *client) newMultipartRequest(ctx context.Context, endpoint *url.URL, body io.Reader, contentType string, contentLength int64) *http.Request {
	req := &http.Request{
		Method: http.MethodPost,
		URL:    endpoint,
		Header: make(http.Header),
		Body:   io.NopCloser(body),
	}
	req.ContentLength = contentLength
	req = req.WithContext(ctx)
	req.Header.Add(xiAPIKeyHeader, c.apikey)
	req.Header.Set("Content-Type", contentType)
	return req
}

func (c *client) doPOSTMultipart(ctx context.Context, endpoint *url.URL, body io.Reader, contentType string, contentLength int64) ([]byte, error) {
	return c.doPOSTMultipartAccepted(ctx, endpoint, body, contentType, contentLength, http.StatusOK)
}

func (c *client) doPOSTMultipartAccepted(ctx context.Context, endpoint *url.URL, body io.Reader, contentType string, contentLength int64, accepted ...int) ([]byte, error) {
	req := c.newMultipartRequest(ctx, endpoint, body, contentType, contentLength)

	clnt := *c.clnt
	clnt.Timeout = 0
	resp, err := clnt.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !statusAccepted(resp.StatusCode, accepted) {
		return nil, fmt.Errorf("api response: %s: %s", resp.Status, respBody)
	}
	return respBody, nil
}

func statusAccepted(code int, accepted []int) bool {
	if len(accepted) == 0 {
		return code == http.StatusOK
	}
	for _, a := range accepted {
		if code == a {
			return true
		}
	}
	return false
}

func (c *client) doPOSTMultipartStream(ctx context.Context, endpoint *url.URL, body io.Reader, contentType string, contentLength int64) (io.ReadCloser, error) {
	req := c.newMultipartRequest(ctx, endpoint, body, contentType, contentLength)
	return c.doStream(req)
}

func (c *client) doStream(req *http.Request) (io.ReadCloser, error) {
	clnt := *c.clnt
	clnt.Timeout = 0
	resp, err := clnt.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api response: %s: %s", resp.Status, body)
	}
	return resp.Body, nil
}

func (c *client) doGET(ctx context.Context, endpoint *url.URL) ([]byte, error) {
	return c.doNoBody(ctx, http.MethodGet, endpoint)
}

func (c *client) doGETNoTimeout(ctx context.Context, endpoint *url.URL) ([]byte, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    endpoint,
		Header: make(http.Header),
	}
	req = req.WithContext(ctx)
	req.Header.Add(xiAPIKeyHeader, c.apikey)
	req.Header.Add("Content-Type", "Application/json")

	clnt := *c.clnt
	clnt.Timeout = 0
	resp, err := clnt.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api response: %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *client) doDELETE(ctx context.Context, endpoint *url.URL) ([]byte, error) {
	return c.doNoBody(ctx, http.MethodDelete, endpoint)
}

func (c *client) doNoBody(ctx context.Context, method string, endpoint *url.URL) ([]byte, error) {
	req := &http.Request{
		Method: method,
		URL:    endpoint,
		Header: make(http.Header),
	}
	req = req.WithContext(ctx)
	req.Header.Add(xiAPIKeyHeader, c.apikey)
	req.Header.Add("Content-Type", "Application/json")

	resp, err := c.clnt.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api response: %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *client) ListModelsWithContext(ctx context.Context) (Models, error) {
	endpoint := c.base.JoinPath(listModelsEndpoint)
	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	models := make(Models, 0)
	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("failed to marshal response to models: %v", err)
	}
	return models, nil
}

func (c *client) ListModels() (Models, error) {
	return c.ListModelsWithContext(context.Background())
}

func (c *client) SearchVoicesWithContext(ctx context.Context, options ...SearchVoicesOption) (VoiceSearchResult, error) {
	var opt searchVoicesOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(searchVoicesEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return VoiceSearchResult{}, err
	}

	var result VoiceSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return VoiceSearchResult{}, fmt.Errorf("failed to marshal response to voices: %v", err)
	}
	return result, nil
}

func (c *client) SearchVoices(options ...SearchVoicesOption) (VoiceSearchResult, error) {
	return c.SearchVoicesWithContext(context.Background(), options...)
}

func (c *client) TextToSpeechWithContext(ctx context.Context, text string, voiceId string, options ...TextToSpeechOption) ([]byte, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	if voiceId == "" {
		return nil, fmt.Errorf("voiceId is empty")
	}

	var opt textToSpeechOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(ttsEndpoint, voiceId)
	endpoint.RawQuery = opt.query().Encode()

	return c.doPOST(ctx, endpoint, opt.body(text))
}

func (c *client) TextToSpeech(text string, voiceId string, options ...TextToSpeechOption) ([]byte, error) {
	return c.TextToSpeechWithContext(context.Background(), text, voiceId, options...)
}

func (c *client) TextToSpeechStreamWithContext(ctx context.Context, text string, voiceId string, options ...TextToSpeechOption) (io.ReadCloser, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	if voiceId == "" {
		return nil, fmt.Errorf("voiceId is empty")
	}

	var opt textToSpeechOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(ttsEndpoint, voiceId, ttsStreamSegment)
	endpoint.RawQuery = opt.query().Encode()

	return c.doPOSTStream(ctx, endpoint, opt.body(text))
}

func (c *client) TextToSpeechStream(text string, voiceId string, options ...TextToSpeechOption) (io.ReadCloser, error) {
	return c.TextToSpeechStreamWithContext(context.Background(), text, voiceId, options...)
}

func (c *client) SpeechToSpeechWithContext(ctx context.Context, audio io.Reader, voiceId string, options ...SpeechToSpeechOption) ([]byte, error) {
	endpoint, body, contentType, err := c.prepareSpeechToSpeech(audio, voiceId, false, options)
	if err != nil {
		return nil, err
	}
	return c.doPOSTMultipart(ctx, endpoint, bytes.NewReader(body), contentType, int64(len(body)))
}

func (c *client) SpeechToSpeech(audio io.Reader, voiceId string, options ...SpeechToSpeechOption) ([]byte, error) {
	return c.SpeechToSpeechWithContext(context.Background(), audio, voiceId, options...)
}

func (c *client) SpeechToSpeechStreamWithContext(ctx context.Context, audio io.Reader, voiceId string, options ...SpeechToSpeechOption) (io.ReadCloser, error) {
	endpoint, body, contentType, err := c.prepareSpeechToSpeech(audio, voiceId, true, options)
	if err != nil {
		return nil, err
	}
	return c.doPOSTMultipartStream(ctx, endpoint, bytes.NewReader(body), contentType, int64(len(body)))
}

func (c *client) SpeechToSpeechStream(audio io.Reader, voiceId string, options ...SpeechToSpeechOption) (io.ReadCloser, error) {
	return c.SpeechToSpeechStreamWithContext(context.Background(), audio, voiceId, options...)
}

func (c *client) prepareSpeechToSpeech(audio io.Reader, voiceId string, stream bool, options []SpeechToSpeechOption) (*url.URL, []byte, string, error) {
	if audio == nil {
		return nil, nil, "", fmt.Errorf("audio is nil")
	}
	if voiceId == "" {
		return nil, nil, "", fmt.Errorf("voiceId is empty")
	}

	var opt speechToSpeechOption
	for _, option := range options {
		option(&opt)
	}

	filename := opt.AudioFilename
	if filename == "" {
		filename = audioFilename(audio)
	}

	body, contentType, err := opt.body(audio, filename)
	if err != nil {
		return nil, nil, "", err
	}

	var endpoint *url.URL
	if stream {
		endpoint = c.base.JoinPath(speechToSpeechEndpoint, voiceId, ttsStreamSegment)
	} else {
		endpoint = c.base.JoinPath(speechToSpeechEndpoint, voiceId)
	}
	endpoint.RawQuery = opt.query().Encode()
	return endpoint, body, contentType, nil
}

func audioFilename(audio io.Reader) string {
	if f, ok := audio.(*os.File); ok {
		name := filepath.Base(f.Name())
		if name != "" && name != "." && name != string(filepath.Separator) {
			return name
		}
	}
	return "audio"
}

func (c *client) SpeechToTextWithContext(ctx context.Context, audio io.Reader, options ...SpeechToTextOption) (SpeechToTextResult, error) {
	endpoint, body, contentType, err := c.prepareSpeechToText(audio, options)
	if err != nil {
		return SpeechToTextResult{}, err
	}
	respBody, err := c.doPOSTMultipartAccepted(ctx, endpoint, bytes.NewReader(body), contentType, int64(len(body)), http.StatusOK, http.StatusAccepted)
	if err != nil {
		return SpeechToTextResult{}, err
	}

	var result SpeechToTextResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return SpeechToTextResult{}, fmt.Errorf("failed to marshal response to transcript: %v", err)
	}
	return result, nil
}

func (c *client) SpeechToText(audio io.Reader, options ...SpeechToTextOption) (SpeechToTextResult, error) {
	return c.SpeechToTextWithContext(context.Background(), audio, options...)
}

func (c *client) prepareSpeechToText(audio io.Reader, options []SpeechToTextOption) (*url.URL, []byte, string, error) {
	var opt speechToTextOption
	for _, option := range options {
		option(&opt)
	}
	if opt.ModelId == "" {
		opt.ModelId = defaultSpeechToTextModel
	}

	hasFile := audio != nil
	hasURL := opt.SourceURL != ""
	if hasFile && hasURL {
		return nil, nil, "", fmt.Errorf("provide audio or source URL, not both")
	}
	if !hasFile && !hasURL {
		return nil, nil, "", fmt.Errorf("audio is nil")
	}

	filename := opt.AudioFilename
	if filename == "" {
		filename = audioFilename(audio)
	}

	body, contentType, err := opt.body(audio, filename)
	if err != nil {
		return nil, nil, "", err
	}

	endpoint := c.base.JoinPath(speechToTextEndpoint)
	endpoint.RawQuery = opt.query().Encode()
	return endpoint, body, contentType, nil
}

func (c *client) GetTranscriptWithContext(ctx context.Context, transcriptionID string) (SpeechToTextResult, error) {
	if transcriptionID == "" {
		return SpeechToTextResult{}, fmt.Errorf("transcriptionID is empty")
	}

	endpoint := c.base.JoinPath(speechToTextTranscriptsPath, transcriptionID)
	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return SpeechToTextResult{}, err
	}

	var result SpeechToTextResult
	if err := json.Unmarshal(body, &result); err != nil {
		return SpeechToTextResult{}, fmt.Errorf("failed to marshal response to transcript: %v", err)
	}
	return result, nil
}

func (c *client) GetTranscript(transcriptionID string) (SpeechToTextResult, error) {
	return c.GetTranscriptWithContext(context.Background(), transcriptionID)
}

func (c *client) DeleteTranscriptWithContext(ctx context.Context, transcriptionID string) error {
	if transcriptionID == "" {
		return fmt.Errorf("transcriptionID is empty")
	}

	endpoint := c.base.JoinPath(speechToTextTranscriptsPath, transcriptionID)
	_, err := c.doDELETE(ctx, endpoint)
	return err
}

func (c *client) DeleteTranscript(transcriptionID string) error {
	return c.DeleteTranscriptWithContext(context.Background(), transcriptionID)
}

func (c *client) ListWebhooksWithContext(ctx context.Context, options ...ListWebhooksOption) (WebhookList, error) {
	var opt listWebhooksOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(workspaceWebhooksEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return WebhookList{}, err
	}

	var result WebhookList
	if err := json.Unmarshal(body, &result); err != nil {
		return WebhookList{}, fmt.Errorf("failed to marshal response to webhooks: %v", err)
	}
	return result, nil
}

func (c *client) ListWebhooks(options ...ListWebhooksOption) (WebhookList, error) {
	return c.ListWebhooksWithContext(context.Background(), options...)
}

func (c *client) CreateWebhookWithContext(ctx context.Context, name string, webhookURL string, options ...CreateWebhookOption) (CreateWebhookResult, error) {
	if name == "" {
		return CreateWebhookResult{}, fmt.Errorf("name is empty")
	}
	if webhookURL == "" {
		return CreateWebhookResult{}, fmt.Errorf("webhookURL is empty")
	}

	var opt createWebhookOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(workspaceWebhooksEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body(name, webhookURL))
	if err != nil {
		return CreateWebhookResult{}, err
	}

	var result CreateWebhookResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateWebhookResult{}, fmt.Errorf("failed to marshal response to webhook: %v", err)
	}
	return result, nil
}

func (c *client) CreateWebhook(name string, webhookURL string, options ...CreateWebhookOption) (CreateWebhookResult, error) {
	return c.CreateWebhookWithContext(context.Background(), name, webhookURL, options...)
}

func (c *client) UpdateWebhookWithContext(ctx context.Context, webhookID string, name string, isDisabled bool, options ...UpdateWebhookOption) (WebhookStatusResult, error) {
	if webhookID == "" {
		return WebhookStatusResult{}, fmt.Errorf("webhookID is empty")
	}
	if name == "" {
		return WebhookStatusResult{}, fmt.Errorf("name is empty")
	}

	var opt updateWebhookOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(workspaceWebhooksEndpoint, webhookID)
	body, err := c.doPATCH(ctx, endpoint, opt.body(name, isDisabled))
	if err != nil {
		return WebhookStatusResult{}, err
	}

	var result WebhookStatusResult
	if err := json.Unmarshal(body, &result); err != nil {
		return WebhookStatusResult{}, fmt.Errorf("failed to marshal response to webhook status: %v", err)
	}
	return result, nil
}

func (c *client) UpdateWebhook(webhookID string, name string, isDisabled bool, options ...UpdateWebhookOption) (WebhookStatusResult, error) {
	return c.UpdateWebhookWithContext(context.Background(), webhookID, name, isDisabled, options...)
}

func (c *client) DeleteWebhookWithContext(ctx context.Context, webhookID string) (WebhookStatusResult, error) {
	if webhookID == "" {
		return WebhookStatusResult{}, fmt.Errorf("webhookID is empty")
	}

	endpoint := c.base.JoinPath(workspaceWebhooksEndpoint, webhookID)
	body, err := c.doDELETE(ctx, endpoint)
	if err != nil {
		return WebhookStatusResult{}, err
	}

	var result WebhookStatusResult
	if err := json.Unmarshal(body, &result); err != nil {
		return WebhookStatusResult{}, fmt.Errorf("failed to marshal response to webhook status: %v", err)
	}
	return result, nil
}

func (c *client) DeleteWebhook(webhookID string) (WebhookStatusResult, error) {
	return c.DeleteWebhookWithContext(context.Background(), webhookID)
}

func (c *client) ListAgentsWithContext(ctx context.Context, options ...ListAgentsOption) (AgentList, error) {
	var opt listAgentsOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(agentsEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return AgentList{}, err
	}

	var result AgentList
	if err := json.Unmarshal(body, &result); err != nil {
		return AgentList{}, fmt.Errorf("failed to marshal response to agents: %v", err)
	}
	return result, nil
}

func (c *client) ListAgents(options ...ListAgentsOption) (AgentList, error) {
	return c.ListAgentsWithContext(context.Background(), options...)
}

func (c *client) CreateAgentWithContext(ctx context.Context, config ConversationalConfig, options ...CreateAgentOption) (CreateAgentResult, error) {
	var opt createAgentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(createAgentEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body(config))
	if err != nil {
		return CreateAgentResult{}, err
	}

	var result CreateAgentResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateAgentResult{}, fmt.Errorf("failed to marshal response to agent: %v", err)
	}
	return result, nil
}

func (c *client) CreateAgent(config ConversationalConfig, options ...CreateAgentOption) (CreateAgentResult, error) {
	return c.CreateAgentWithContext(context.Background(), config, options...)
}

func (c *client) CreateOrUpdateAgentWithContext(ctx context.Context, config ConversationalConfig, options ...CreateAgentOption) (CreateOrUpdateAgentResult, error) {
	var opt createAgentOption
	for _, option := range options {
		option(&opt)
	}
	if opt.Name == "" {
		return CreateOrUpdateAgentResult{}, fmt.Errorf("agent name is empty")
	}

	matches, err := c.findAgentsByExactName(ctx, opt.Name)
	if err != nil {
		return CreateOrUpdateAgentResult{}, err
	}
	if len(matches) > 1 {
		return CreateOrUpdateAgentResult{}, fmt.Errorf("multiple agents named %q", opt.Name)
	}
	if len(matches) == 0 {
		created, err := c.CreateAgentWithContext(ctx, config, options...)
		if err != nil {
			return CreateOrUpdateAgentResult{}, err
		}
		return CreateOrUpdateAgentResult{AgentID: created.AgentID, Created: true}, nil
	}

	updated, err := c.UpdateAgentWithContext(ctx, matches[0].AgentID, opt.updateOptions(config)...)
	if err != nil {
		return CreateOrUpdateAgentResult{}, err
	}
	return CreateOrUpdateAgentResult{AgentID: updated.AgentID}, nil
}

func (c *client) CreateOrUpdateAgent(config ConversationalConfig, options ...CreateAgentOption) (CreateOrUpdateAgentResult, error) {
	return c.CreateOrUpdateAgentWithContext(context.Background(), config, options...)
}

func (c *client) findAgentsByExactName(ctx context.Context, name string) ([]AgentSummary, error) {
	var matches []AgentSummary
	var cursor string
	for {
		opts := []ListAgentsOption{WithAgentSearch(name), WithAgentPageSize(100)}
		if cursor != "" {
			opts = append(opts, WithAgentCursor(cursor))
		}

		page, err := c.ListAgentsWithContext(ctx, opts...)
		if err != nil {
			return nil, err
		}
		for _, agent := range page.Agents {
			if agent.Name == name {
				matches = append(matches, agent)
			}
		}
		if !page.HasMore || page.NextCursor == "" || page.NextCursor == cursor {
			return matches, nil
		}
		cursor = page.NextCursor
	}
}

func (c *client) GetAgentWithContext(ctx context.Context, agentID string, options ...GetAgentOption) (Agent, error) {
	if agentID == "" {
		return Agent{}, fmt.Errorf("agentID is empty")
	}

	var opt getAgentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(agentsEndpoint, agentID)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return Agent{}, err
	}

	var result Agent
	if err := json.Unmarshal(body, &result); err != nil {
		return Agent{}, fmt.Errorf("failed to marshal response to agent: %v", err)
	}
	return result, nil
}

func (c *client) GetAgent(agentID string, options ...GetAgentOption) (Agent, error) {
	return c.GetAgentWithContext(context.Background(), agentID, options...)
}

func (c *client) UpdateAgentWithContext(ctx context.Context, agentID string, options ...UpdateAgentOption) (Agent, error) {
	if agentID == "" {
		return Agent{}, fmt.Errorf("agentID is empty")
	}

	var opt updateAgentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(agentsEndpoint, agentID)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doPATCH(ctx, endpoint, opt.body())
	if err != nil {
		return Agent{}, err
	}

	var result Agent
	if err := json.Unmarshal(body, &result); err != nil {
		return Agent{}, fmt.Errorf("failed to marshal response to agent: %v", err)
	}
	return result, nil
}

func (c *client) UpdateAgent(agentID string, options ...UpdateAgentOption) (Agent, error) {
	return c.UpdateAgentWithContext(context.Background(), agentID, options...)
}

func (c *client) DeleteAgentWithContext(ctx context.Context, agentID string, options ...DeleteAgentOption) error {
	if agentID == "" {
		return fmt.Errorf("agentID is empty")
	}

	var opt deleteAgentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(agentsEndpoint, agentID)
	endpoint.RawQuery = opt.query().Encode()
	_, err := c.doDELETE(ctx, endpoint)
	return err
}

func (c *client) DeleteAgent(agentID string, options ...DeleteAgentOption) error {
	return c.DeleteAgentWithContext(context.Background(), agentID, options...)
}

func (c *client) DuplicateAgentWithContext(ctx context.Context, agentID string, options ...DuplicateAgentOption) (CreateAgentResult, error) {
	if agentID == "" {
		return CreateAgentResult{}, fmt.Errorf("agentID is empty")
	}

	var opt duplicateAgentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(agentsEndpoint, agentID, "duplicate")
	body, err := c.doPOST(ctx, endpoint, opt.body())
	if err != nil {
		return CreateAgentResult{}, err
	}

	var result CreateAgentResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateAgentResult{}, fmt.Errorf("failed to marshal response to agent: %v", err)
	}
	return result, nil
}

func (c *client) DuplicateAgent(agentID string, options ...DuplicateAgentOption) (CreateAgentResult, error) {
	return c.DuplicateAgentWithContext(context.Background(), agentID, options...)
}

func (c *client) GetAgentLinkWithContext(ctx context.Context, agentID string) (AgentLink, error) {
	if agentID == "" {
		return AgentLink{}, fmt.Errorf("agentID is empty")
	}

	endpoint := c.base.JoinPath(agentsEndpoint, agentID, "link")
	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return AgentLink{}, err
	}

	var result AgentLink
	if err := json.Unmarshal(body, &result); err != nil {
		return AgentLink{}, fmt.Errorf("failed to marshal response to agent link: %v", err)
	}
	return result, nil
}

func (c *client) GetAgentLink(agentID string) (AgentLink, error) {
	return c.GetAgentLinkWithContext(context.Background(), agentID)
}

func (c *client) ListConversationsWithContext(ctx context.Context, options ...ListConversationsOption) (ConversationList, error) {
	var opt listConversationsOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(conversationsEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return ConversationList{}, err
	}

	var result ConversationList
	if err := json.Unmarshal(body, &result); err != nil {
		return ConversationList{}, fmt.Errorf("failed to marshal response to conversations: %v", err)
	}
	return result, nil
}

func (c *client) ListConversations(options ...ListConversationsOption) (ConversationList, error) {
	return c.ListConversationsWithContext(context.Background(), options...)
}

func (c *client) GetConversationWithContext(ctx context.Context, conversationID string) (Conversation, error) {
	if conversationID == "" {
		return Conversation{}, fmt.Errorf("conversationID is empty")
	}

	endpoint := c.base.JoinPath(conversationsEndpoint, conversationID)
	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return Conversation{}, err
	}

	var result Conversation
	if err := json.Unmarshal(body, &result); err != nil {
		return Conversation{}, fmt.Errorf("failed to marshal response to conversation: %v", err)
	}
	return result, nil
}

func (c *client) GetConversation(conversationID string) (Conversation, error) {
	return c.GetConversationWithContext(context.Background(), conversationID)
}

func (c *client) GetConversationSummaryWithContext(ctx context.Context, conversationID string, options ...GetConversationSummaryOption) (ConversationSummary, error) {
	if conversationID == "" {
		return ConversationSummary{}, fmt.Errorf("conversationID is empty")
	}

	var opt getConversationSummaryOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(conversationsEndpoint, conversationID, "summary")
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return ConversationSummary{}, err
	}

	var result ConversationSummary
	if err := json.Unmarshal(body, &result); err != nil {
		return ConversationSummary{}, fmt.Errorf("failed to marshal response to conversation summary: %v", err)
	}
	return result, nil
}

func (c *client) GetConversationSummary(conversationID string, options ...GetConversationSummaryOption) (ConversationSummary, error) {
	return c.GetConversationSummaryWithContext(context.Background(), conversationID, options...)
}

func (c *client) DeleteConversationWithContext(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return fmt.Errorf("conversationID is empty")
	}

	endpoint := c.base.JoinPath(conversationsEndpoint, conversationID)
	_, err := c.doDELETE(ctx, endpoint)
	return err
}

func (c *client) DeleteConversation(conversationID string) error {
	return c.DeleteConversationWithContext(context.Background(), conversationID)
}

func (c *client) GetConversationAudioWithContext(ctx context.Context, conversationID string) ([]byte, error) {
	if conversationID == "" {
		return nil, fmt.Errorf("conversationID is empty")
	}

	endpoint := c.base.JoinPath(conversationsEndpoint, conversationID, "audio")
	return c.doGETNoTimeout(ctx, endpoint)
}

func (c *client) GetConversationAudio(conversationID string) ([]byte, error) {
	return c.GetConversationAudioWithContext(context.Background(), conversationID)
}

func (c *client) GetSignedURLWithContext(ctx context.Context, agentID string, options ...SignedURLOption) (SignedURL, error) {
	if agentID == "" {
		return SignedURL{}, fmt.Errorf("agentID is empty")
	}

	var opt signedURLOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(conversationSignedURLEndpoint)
	endpoint.RawQuery = opt.query(agentID).Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return SignedURL{}, err
	}

	var result SignedURL
	if err := json.Unmarshal(body, &result); err != nil {
		return SignedURL{}, fmt.Errorf("failed to marshal response to signed URL: %v", err)
	}
	return result, nil
}

func (c *client) GetSignedURL(agentID string, options ...SignedURLOption) (SignedURL, error) {
	return c.GetSignedURLWithContext(context.Background(), agentID, options...)
}

func (c *client) GetConversationTokenWithContext(ctx context.Context, agentID string, options ...ConversationTokenOption) (ConversationToken, error) {
	if agentID == "" {
		return ConversationToken{}, fmt.Errorf("agentID is empty")
	}

	var opt conversationTokenOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(conversationTokenEndpoint)
	endpoint.RawQuery = opt.query(agentID).Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return ConversationToken{}, err
	}

	var result ConversationToken
	if err := json.Unmarshal(body, &result); err != nil {
		return ConversationToken{}, fmt.Errorf("failed to marshal response to conversation token: %v", err)
	}
	return result, nil
}

func (c *client) GetConversationToken(agentID string, options ...ConversationTokenOption) (ConversationToken, error) {
	return c.GetConversationTokenWithContext(context.Background(), agentID, options...)
}

func (c *client) SendConversationFeedbackWithContext(ctx context.Context, conversationID string, feedback ConversationFeedback) error {
	if conversationID == "" {
		return fmt.Errorf("conversationID is empty")
	}
	if feedback == "" {
		return fmt.Errorf("feedback is empty")
	}

	endpoint := c.base.JoinPath(conversationsEndpoint, conversationID, "feedback")
	_, err := c.doPOST(ctx, endpoint, conversationFeedbackRequest{Feedback: feedback.String()})
	return err
}

func (c *client) SendConversationFeedback(conversationID string, feedback ConversationFeedback) error {
	return c.SendConversationFeedbackWithContext(context.Background(), conversationID, feedback)
}

func (c *client) CreateRealtimeClientSecretWithContext(ctx context.Context, options ...ClientSecretOption) (RealtimeClientSecret, error) {
	var opt clientSecretOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(realtimeClientSecretsEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body())
	if err != nil {
		return RealtimeClientSecret{}, err
	}

	var result RealtimeClientSecret
	if err := json.Unmarshal(body, &result); err != nil {
		return RealtimeClientSecret{}, fmt.Errorf("failed to marshal response to client secret: %v", err)
	}
	return result, nil
}

func (c *client) CreateRealtimeClientSecret(options ...ClientSecretOption) (RealtimeClientSecret, error) {
	return c.CreateRealtimeClientSecretWithContext(context.Background(), options...)
}

func (c *client) ListToolsWithContext(ctx context.Context, options ...ListToolsOption) (ToolList, error) {
	var opt listToolsOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return ToolList{}, err
	}

	var result ToolList
	if err := json.Unmarshal(body, &result); err != nil {
		return ToolList{}, fmt.Errorf("failed to marshal response to tools: %v", err)
	}
	return result, nil
}

func (c *client) ListTools(options ...ListToolsOption) (ToolList, error) {
	return c.ListToolsWithContext(context.Background(), options...)
}

func (c *client) CreateToolWithContext(ctx context.Context, config ToolConfig, options ...CreateToolOption) (Tool, error) {
	var opt createToolOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body(config))
	if err != nil {
		return Tool{}, err
	}

	var result Tool
	if err := json.Unmarshal(body, &result); err != nil {
		return Tool{}, fmt.Errorf("failed to marshal response to tool: %v", err)
	}
	return result, nil
}

func (c *client) CreateTool(config ToolConfig, options ...CreateToolOption) (Tool, error) {
	return c.CreateToolWithContext(context.Background(), config, options...)
}

func (c *client) GetToolWithContext(ctx context.Context, toolID string, options ...GetToolOption) (Tool, error) {
	if toolID == "" {
		return Tool{}, fmt.Errorf("toolID is empty")
	}

	var opt getToolOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint, toolID)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return Tool{}, err
	}

	var result Tool
	if err := json.Unmarshal(body, &result); err != nil {
		return Tool{}, fmt.Errorf("failed to marshal response to tool: %v", err)
	}
	return result, nil
}

func (c *client) GetTool(toolID string, options ...GetToolOption) (Tool, error) {
	return c.GetToolWithContext(context.Background(), toolID, options...)
}

func (c *client) UpdateToolWithContext(ctx context.Context, toolID string, config ToolConfig, options ...UpdateToolOption) (Tool, error) {
	if toolID == "" {
		return Tool{}, fmt.Errorf("toolID is empty")
	}

	var opt createToolOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint, toolID)
	body, err := c.doPATCH(ctx, endpoint, opt.body(config))
	if err != nil {
		return Tool{}, err
	}

	var result Tool
	if err := json.Unmarshal(body, &result); err != nil {
		return Tool{}, fmt.Errorf("failed to marshal response to tool: %v", err)
	}
	return result, nil
}

func (c *client) UpdateTool(toolID string, config ToolConfig, options ...UpdateToolOption) (Tool, error) {
	return c.UpdateToolWithContext(context.Background(), toolID, config, options...)
}

func (c *client) DeleteToolWithContext(ctx context.Context, toolID string, options ...DeleteToolOption) error {
	if toolID == "" {
		return fmt.Errorf("toolID is empty")
	}

	var opt deleteToolOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint, toolID)
	endpoint.RawQuery = opt.query().Encode()
	_, err := c.doDELETE(ctx, endpoint)
	return err
}

func (c *client) DeleteTool(toolID string, options ...DeleteToolOption) error {
	return c.DeleteToolWithContext(context.Background(), toolID, options...)
}

func (c *client) GetToolDependentAgentsWithContext(ctx context.Context, toolID string, options ...ToolDependentOption) (DependentAgentList, error) {
	if toolID == "" {
		return DependentAgentList{}, fmt.Errorf("toolID is empty")
	}

	var opt toolDependentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(toolsEndpoint, toolID, "dependent-agents")
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return DependentAgentList{}, err
	}

	var result DependentAgentList
	if err := json.Unmarshal(body, &result); err != nil {
		return DependentAgentList{}, fmt.Errorf("failed to marshal response to dependent agents: %v", err)
	}
	return result, nil
}

func (c *client) GetToolDependentAgents(toolID string, options ...ToolDependentOption) (DependentAgentList, error) {
	return c.GetToolDependentAgentsWithContext(context.Background(), toolID, options...)
}

func (c *client) ListKnowledgeBaseWithContext(ctx context.Context, options ...ListKnowledgeOption) (KnowledgeList, error) {
	var opt listKnowledgeOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return KnowledgeList{}, err
	}

	var result KnowledgeList
	if err := json.Unmarshal(body, &result); err != nil {
		return KnowledgeList{}, fmt.Errorf("failed to marshal response to knowledge base: %v", err)
	}
	return result, nil
}

func (c *client) ListKnowledgeBase(options ...ListKnowledgeOption) (KnowledgeList, error) {
	return c.ListKnowledgeBaseWithContext(context.Background(), options...)
}

func (c *client) CreateKnowledgeBaseTextWithContext(ctx context.Context, text string, options ...CreateKnowledgeTextOption) (CreateKnowledgeResult, error) {
	if text == "" {
		return CreateKnowledgeResult{}, fmt.Errorf("text is empty")
	}

	var opt createKnowledgeTextOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseTextEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body(text))
	if err != nil {
		return CreateKnowledgeResult{}, err
	}

	var result CreateKnowledgeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateKnowledgeResult{}, fmt.Errorf("failed to marshal response to knowledge base: %v", err)
	}
	return result, nil
}

func (c *client) CreateKnowledgeBaseText(text string, options ...CreateKnowledgeTextOption) (CreateKnowledgeResult, error) {
	return c.CreateKnowledgeBaseTextWithContext(context.Background(), text, options...)
}

func (c *client) CreateKnowledgeBaseURLWithContext(ctx context.Context, pageURL string, options ...CreateKnowledgeURLOption) (CreateKnowledgeResult, error) {
	if pageURL == "" {
		return CreateKnowledgeResult{}, fmt.Errorf("url is empty")
	}

	var opt createKnowledgeURLOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseURLEndpoint)
	body, err := c.doPOST(ctx, endpoint, opt.body(pageURL))
	if err != nil {
		return CreateKnowledgeResult{}, err
	}

	var result CreateKnowledgeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateKnowledgeResult{}, fmt.Errorf("failed to marshal response to knowledge base: %v", err)
	}
	return result, nil
}

func (c *client) CreateKnowledgeBaseURL(pageURL string, options ...CreateKnowledgeURLOption) (CreateKnowledgeResult, error) {
	return c.CreateKnowledgeBaseURLWithContext(context.Background(), pageURL, options...)
}

func (c *client) CreateKnowledgeBaseFileWithContext(ctx context.Context, file io.Reader, options ...CreateKnowledgeFileOption) (CreateKnowledgeResult, error) {
	if file == nil {
		return CreateKnowledgeResult{}, fmt.Errorf("file is nil")
	}

	var opt createKnowledgeFileOption
	for _, option := range options {
		option(&opt)
	}

	filename := opt.Filename
	if filename == "" {
		filename = formFilename(file, "document")
	}

	payload, contentType, err := opt.body(file, filename)
	if err != nil {
		return CreateKnowledgeResult{}, err
	}

	endpoint := c.base.JoinPath(knowledgeBaseFileEndpoint)
	body, err := c.doPOSTMultipart(ctx, endpoint, bytes.NewReader(payload), contentType, int64(len(payload)))
	if err != nil {
		return CreateKnowledgeResult{}, err
	}

	var result CreateKnowledgeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return CreateKnowledgeResult{}, fmt.Errorf("failed to marshal response to knowledge base: %v", err)
	}
	return result, nil
}

func (c *client) CreateKnowledgeBaseFile(file io.Reader, options ...CreateKnowledgeFileOption) (CreateKnowledgeResult, error) {
	return c.CreateKnowledgeBaseFileWithContext(context.Background(), file, options...)
}

func formFilename(file io.Reader, fallback string) string {
	if f, ok := file.(*os.File); ok {
		name := filepath.Base(f.Name())
		if name != "" && name != "." && name != string(filepath.Separator) {
			return name
		}
	}
	return fallback
}

func (c *client) GetKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...GetKnowledgeOption) (KnowledgeDocument, error) {
	if documentationID == "" {
		return KnowledgeDocument{}, fmt.Errorf("documentationID is empty")
	}

	var opt getKnowledgeOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return KnowledgeDocument{}, err
	}

	var result KnowledgeDocument
	if err := json.Unmarshal(body, &result); err != nil {
		return KnowledgeDocument{}, fmt.Errorf("failed to marshal response to knowledge base document: %v", err)
	}
	return result, nil
}

func (c *client) GetKnowledgeBase(documentationID string, options ...GetKnowledgeOption) (KnowledgeDocument, error) {
	return c.GetKnowledgeBaseWithContext(context.Background(), documentationID, options...)
}

func (c *client) UpdateKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...UpdateKnowledgeOption) (KnowledgeDocument, error) {
	if documentationID == "" {
		return KnowledgeDocument{}, fmt.Errorf("documentationID is empty")
	}

	var opt updateKnowledgeOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID)
	body, err := c.doPATCH(ctx, endpoint, opt.body())
	if err != nil {
		return KnowledgeDocument{}, err
	}

	var result KnowledgeDocument
	if err := json.Unmarshal(body, &result); err != nil {
		return KnowledgeDocument{}, fmt.Errorf("failed to marshal response to knowledge base document: %v", err)
	}
	return result, nil
}

func (c *client) UpdateKnowledgeBase(documentationID string, options ...UpdateKnowledgeOption) (KnowledgeDocument, error) {
	return c.UpdateKnowledgeBaseWithContext(context.Background(), documentationID, options...)
}

func (c *client) DeleteKnowledgeBaseWithContext(ctx context.Context, documentationID string, options ...DeleteKnowledgeOption) error {
	if documentationID == "" {
		return fmt.Errorf("documentationID is empty")
	}

	var opt deleteKnowledgeOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID)
	endpoint.RawQuery = opt.query().Encode()
	_, err := c.doDELETE(ctx, endpoint)
	return err
}

func (c *client) DeleteKnowledgeBase(documentationID string, options ...DeleteKnowledgeOption) error {
	return c.DeleteKnowledgeBaseWithContext(context.Background(), documentationID, options...)
}

func (c *client) GetKnowledgeBaseContentWithContext(ctx context.Context, documentationID string) (KnowledgeContent, error) {
	if documentationID == "" {
		return KnowledgeContent{}, fmt.Errorf("documentationID is empty")
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID, "content")
	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return KnowledgeContent{}, err
	}
	return decodeKnowledgeContent(body), nil
}

func (c *client) GetKnowledgeBaseContent(documentationID string) (KnowledgeContent, error) {
	return c.GetKnowledgeBaseContentWithContext(context.Background(), documentationID)
}

func (c *client) GetKnowledgeBaseDependentAgentsWithContext(ctx context.Context, documentationID string, options ...KnowledgeDependentOption) (DependentAgentList, error) {
	if documentationID == "" {
		return DependentAgentList{}, fmt.Errorf("documentationID is empty")
	}

	var opt knowledgeDependentOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID, "dependent-agents")
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return DependentAgentList{}, err
	}

	var result DependentAgentList
	if err := json.Unmarshal(body, &result); err != nil {
		return DependentAgentList{}, fmt.Errorf("failed to marshal response to dependent agents: %v", err)
	}
	return result, nil
}

func (c *client) GetKnowledgeBaseDependentAgents(documentationID string, options ...KnowledgeDependentOption) (DependentAgentList, error) {
	return c.GetKnowledgeBaseDependentAgentsWithContext(context.Background(), documentationID, options...)
}

func (c *client) ComputeRAGIndexWithContext(ctx context.Context, documentationID string, model EmbeddingModel) (RAGIndex, error) {
	if documentationID == "" {
		return RAGIndex{}, fmt.Errorf("documentationID is empty")
	}
	if model == "" {
		return RAGIndex{}, fmt.Errorf("model is empty")
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID, "rag-index")
	body, err := c.doPOST(ctx, endpoint, ragIndexRequest{Model: model.String()})
	if err != nil {
		return RAGIndex{}, err
	}

	var result RAGIndex
	if err := json.Unmarshal(body, &result); err != nil {
		return RAGIndex{}, fmt.Errorf("failed to marshal response to rag index: %v", err)
	}
	return result, nil
}

func (c *client) ComputeRAGIndex(documentationID string, model EmbeddingModel) (RAGIndex, error) {
	return c.ComputeRAGIndexWithContext(context.Background(), documentationID, model)
}

func (c *client) GetKnowledgeBaseChunkWithContext(ctx context.Context, documentationID string, chunkID string, options ...GetKnowledgeChunkOption) (KnowledgeChunk, error) {
	if documentationID == "" {
		return KnowledgeChunk{}, fmt.Errorf("documentationID is empty")
	}
	if chunkID == "" {
		return KnowledgeChunk{}, fmt.Errorf("chunkID is empty")
	}

	var opt getKnowledgeChunkOption
	for _, option := range options {
		option(&opt)
	}

	endpoint := c.base.JoinPath(knowledgeBaseEndpoint, documentationID, "chunk", chunkID)
	endpoint.RawQuery = opt.query().Encode()

	body, err := c.doGET(ctx, endpoint)
	if err != nil {
		return KnowledgeChunk{}, err
	}

	var result KnowledgeChunk
	if err := json.Unmarshal(body, &result); err != nil {
		return KnowledgeChunk{}, fmt.Errorf("failed to marshal response to knowledge chunk: %v", err)
	}
	return result, nil
}

func (c *client) GetKnowledgeBaseChunk(documentationID string, chunkID string, options ...GetKnowledgeChunkOption) (KnowledgeChunk, error) {
	return c.GetKnowledgeBaseChunkWithContext(context.Background(), documentationID, chunkID, options...)
}
