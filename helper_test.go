package eleven

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func testClient(t *testing.T, handler http.HandlerFunc) *client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := NewClient("test-key").(*client)
	base, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	c.base = base
	c.clnt = server.Client()
	return c
}

func liveClient(t *testing.T) Client {
	t.Helper()
	if os.Getenv("ELEVEN_LIVE") == "" {
		t.Skip("set ELEVEN_LIVE=1 to hit the API")
	}
	apiKey := os.Getenv("ELEVEN_API_KEY")
	if apiKey == "" {
		t.Fatal("ELEVEN_API_KEY is empty")
	}
	return NewClient(apiKey)
}

func requireHeader(t *testing.T, r *http.Request, key, want string) {
	t.Helper()
	if got := r.Header.Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
