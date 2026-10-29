package eleven

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestOutputFormatString(t *testing.T) {
	cases := []struct {
		format OutputFormat
		want   string
	}{
		{ALAW8000, "alaw_8000"},
		{MP3_22050_32, "mp3_22050_32"},
		{MP3_24000_48, "mp3_24000_48"},
		{MP3_44100_32, "mp3_44100_32"},
		{MP3_44100_64, "mp3_44100_64"},
		{MP3_44100_96, "mp3_44100_96"},
		{MP3_44100_128, "mp3_44100_128"},
		{MP3_44100_192, "mp3_44100_192"},
		{OPUS_48000_32, "opus_48000_32"},
		{OPUS_48000_64, "opus_48000_64"},
		{OPUS_48000_96, "opus_48000_96"},
		{OPUS_48000_128, "opus_48000_128"},
		{OPUS_48000_192, "opus_48000_192"},
		{PCM_8000, "pcm_8000"},
		{PCM_16000, "pcm_16000"},
		{PCM_22050, "pcm_22050"},
		{PCM_24000, "pcm_24000"},
		{PCM_32000, "pcm_32000"},
		{PCM_44100, "pcm_44100"},
		{PCM_48000, "pcm_48000"},
		{ULAW8000, "ulaw_8000"},
		{WAV_8000, "wav_8000"},
		{WAV_16000, "wav_16000"},
		{WAV_22050, "wav_22050"},
		{WAV_24000, "wav_24000"},
		{WAV_32000, "wav_32000"},
		{WAV_44100, "wav_44100"},
		{WAV_48000, "wav_48000"},
		{OutputFormat(255), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.format.String(); got != tc.want {
			t.Fatalf("%d = %q, want %q", tc.format, got, tc.want)
		}
	}
}

func TestStreamingLatencyString(t *testing.T) {
	if StreamingLatencyDefault.String() != "0" {
		t.Fatalf("default = %s", StreamingLatencyDefault)
	}
	if StreamingLatencyMaxNoNormalizer.String() != "4" {
		t.Fatalf("max no normalizer = %s", StreamingLatencyMaxNoNormalizer)
	}
}

func TestTextToSpeechQuery(t *testing.T) {
	var opt textToSpeechOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithoutEnableLogging()(&opt)
	WithOptimizeStreamingLatency(StreamingLatencyStrong)(&opt)
	WithOutputFormat(MP3_44100_128)(&opt)

	q := opt.query()
	if q.Get("enable_logging") != "false" {
		t.Fatalf("enable_logging = %q", q.Get("enable_logging"))
	}
	if q.Get("optimize_streaming_latency") != "2" {
		t.Fatalf("optimize_streaming_latency = %q", q.Get("optimize_streaming_latency"))
	}
	if q.Get("output_format") != "mp3_44100_128" {
		t.Fatalf("output_format = %q", q.Get("output_format"))
	}
}

func TestTextToSpeechRequestBody(t *testing.T) {
	var opt textToSpeechOption
	got, err := json.Marshal(opt.body("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"text":"hello"}` {
		t.Fatalf("minimal body = %s", got)
	}

	WithModelId("eleven_multilingual_v2")(&opt)
	WithLanguageCode("en")(&opt)
	WithVoiceSettings(VoiceSettings{Stability: 0.5, SimilarityBoost: 0.75, Speed: 1})(&opt)
	WithPronunciationDictionaryLocators(PronunciationDictionaryLocator{PronunciationDictionaryId: "dict", VersionId: "v1"})(&opt)
	WithSeed(42)(&opt)
	WithPreviousText("before")(&opt)
	WithNextText("after")(&opt)
	WithPreviousRequestIds("p1")(&opt)
	WithNextRequestIds("n1")(&opt)
	WithUsePvcAsIvc(true)(&opt)
	WithApplyTextNormalization(TextNormalizationOn)(&opt)
	WithApplyLanguageTextNormalization()(&opt)
	WithAudioEffects(AudioEffects{Pan: 0.1})(&opt)

	got, err = json.Marshal(opt.body("hello"))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["text"] != "hello" || payload["model_id"] != "eleven_multilingual_v2" {
		t.Fatalf("payload = %s", got)
	}
	if payload["apply_text_normalization"] != "on" {
		t.Fatalf("normalization = %v", payload["apply_text_normalization"])
	}
	if payload["apply_language_text_normalization"] != true {
		t.Fatalf("language normalization = %v", payload["apply_language_text_normalization"])
	}
	settings := payload["voice_settings"].(map[string]any)
	if settings["stability"] != 0.5 || settings["speed"] != 1.0 {
		t.Fatalf("voice_settings = %v", settings)
	}
	if _, ok := settings["use_speaker_boost"]; !ok {
		t.Fatal("voice_settings omitted zero bool")
	}
}

func TestTextToSpeechValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.TextToSpeech("", "voice"); err == nil || err.Error() != "text is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.TextToSpeech("hello", ""); err == nil || err.Error() != "voiceId is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.TextToSpeechStream("", "voice"); err == nil || err.Error() != "text is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.TextToSpeechStream("hello", ""); err == nil || err.Error() != "voiceId is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestTextToSpeech(t *testing.T) {
	audio := []byte("ID3fake-mp3")
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/text-to-speech/voice-1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		requireHeader(t, r, "Content-Type", "Application/json")
		if r.URL.Query().Get("output_format") != "mp3_44100_128" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		var payload textToSpeechRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Text != "Hello from go eleven." || payload.ModelId != "eleven_multilingual_v2" {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write(audio)
	})

	got, err := client.TextToSpeech(
		"Hello from go eleven.",
		"voice-1",
		WithModelId("eleven_multilingual_v2"),
		WithOutputFormat(MP3_44100_128),
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(audio) {
		t.Fatalf("audio = %q", got)
	}
}

func TestTextToSpeechError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad voice", http.StatusBadRequest)
	})
	_, err := client.TextToSpeech("hello", "voice-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTextToSpeechStream(t *testing.T) {
	audio := []byte("streamed-mp3-bytes")
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/text-to-speech/voice-1/stream" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write(audio)
	})

	stream, err := client.TextToSpeechStream("Hello from go eleven.", "voice-1", WithOutputFormat(MP3_44100_128))
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	got, err := io.ReadAll(stream)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(audio) {
		t.Fatalf("audio = %q", got)
	}
}

func TestTextToSpeechStreamError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "wav not supported", http.StatusBadRequest)
	})
	_, err := client.TextToSpeechStream("hello", "voice-1", WithOutputFormat(WAV_44100))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLiveTextToSpeech(t *testing.T) {
	audio, err := liveClient(t).TextToSpeech(
		"Hello from go eleven.",
		"CwhRBWXzGAHq8TQ4Fs17",
		WithModelId("eleven_multilingual_v2"),
		WithOutputFormat(MP3_44100_128),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(audio) == 0 {
		t.Fatal("empty audio")
	}
	t.Logf("bytes=%d header=%q", len(audio), audio[:min(3, len(audio))])
}

func TestLiveTextToSpeechStream(t *testing.T) {
	stream, err := liveClient(t).TextToSpeechStream(
		"Hello from go eleven.",
		"CwhRBWXzGAHq8TQ4Fs17",
		WithModelId("eleven_multilingual_v2"),
		WithOutputFormat(MP3_44100_128),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	audio, err := io.ReadAll(stream)
	if err != nil {
		t.Fatal(err)
	}
	if len(audio) == 0 {
		t.Fatal("empty stream")
	}
	t.Logf("bytes=%d", len(audio))
}
