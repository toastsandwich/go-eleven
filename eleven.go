package eleven

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	req, err := c.newJSONRequest(ctx, http.MethodPost, endpoint, payload)
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
	req := &http.Request{
		Method: http.MethodGet,
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
