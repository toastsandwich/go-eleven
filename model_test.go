package eleven

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestModelsAll(t *testing.T) {
	models := Models{
		{Name: "a", CanDoTextToSpeech: true},
		{Name: "b", CanDoTextToSpeech: false},
	}

	var names []string
	for model := range models.All() {
		names = append(names, model.Name)
	}
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("names = %v", names)
	}
}

func TestModelsAllWith(t *testing.T) {
	models := Models{
		{Name: "tts", CanDoTextToSpeech: true, CanDoVoiceConversion: false},
		{Name: "vc", CanDoTextToSpeech: false, CanDoVoiceConversion: true},
		{Name: "both", CanDoTextToSpeech: true, CanDoVoiceConversion: true},
	}

	t.Run("zero filters yields all", func(t *testing.T) {
		var names []string
		for model := range models.AllWith() {
			names = append(names, model.Name)
		}
		if len(names) != 3 {
			t.Fatalf("names = %v", names)
		}
	})

	t.Run("all filters must match", func(t *testing.T) {
		var names []string
		for model := range models.AllWith(
			func(m *Model) bool { return m.CanDoTextToSpeech },
			func(m *Model) bool { return m.CanDoVoiceConversion },
		) {
			names = append(names, model.Name)
		}
		if len(names) != 1 || names[0] != "both" {
			t.Fatalf("names = %v", names)
		}
	})

	t.Run("early stop", func(t *testing.T) {
		var names []string
		for model := range models.All() {
			names = append(names, model.Name)
			break
		}
		if len(names) != 1 || names[0] != "tts" {
			t.Fatalf("names = %v", names)
		}
	})
}

func TestListModels(t *testing.T) {
	const body = `[{"model_id":"eleven_multilingual_v2","name":"Multilingual v2","can_do_text_to_speech":true,"languages":[{"language_id":"en","name":"English"}]}]`
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	models, err := client.ListModels()
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 1 {
		t.Fatalf("len = %d", len(models))
	}
	if models[0].ModelId != "eleven_multilingual_v2" || !models[0].CanDoTextToSpeech {
		t.Fatalf("model = %+v", models[0])
	}
	if len(models[0].Languages) != 1 || models[0].Languages[0].LanguageId != "en" {
		t.Fatalf("languages = %+v", models[0].Languages)
	}
}

func TestListModelsError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"detail":"unauthorized"}`, http.StatusUnauthorized)
	})
	_, err := client.ListModels()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestModelJSONRoundTrip(t *testing.T) {
	raw := []byte(`{"model_id":"eleven_turbo_v2","name":"Turbo","can_do_text_to_speech":true,"token_cost_factor":1.5}`)
	var model Model
	if err := json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	if model.ModelId != "eleven_turbo_v2" || model.TokenCostFactor != 1.5 {
		t.Fatalf("model = %+v", model)
	}
}

func TestLiveListModels(t *testing.T) {
	models, err := liveClient(t).ListModels()
	if err != nil {
		t.Fatal(err)
	}
	if len(models) == 0 {
		t.Fatal("no models")
	}
	var tts int
	for model := range models.AllWith(func(m *Model) bool { return m.CanDoTextToSpeech }) {
		tts++
		if model.ModelId == "" || model.Name == "" {
			t.Fatalf("incomplete model: %+v", model)
		}
	}
	if tts == 0 {
		t.Fatal("no text-to-speech models")
	}
	t.Logf("models=%d tts=%d", len(models), tts)
}
