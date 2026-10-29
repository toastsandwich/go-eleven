package eleven

import "testing"

func TestNewClient(t *testing.T) {
	c := NewClient("sk_test")
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	got := c.(*client)
	if got.apikey != "sk_test" {
		t.Fatalf("apikey = %q", got.apikey)
	}
	if got.base.String() != parentEndpoint {
		t.Fatalf("base = %q", got.base)
	}
}

func TestNewClientEmptyAPIKeyPanics(t *testing.T) {
	defer func() {
		if recover() != "apikey is empty" {
			t.Fatal("expected panic")
		}
	}()
	NewClient("")
}
