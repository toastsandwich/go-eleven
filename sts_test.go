package eleven

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileFormatString(t *testing.T) {
	if FileFormatOther.String() != "other" {
		t.Fatalf("other = %s", FileFormatOther)
	}
	if FileFormatPCMS16LE16.String() != "pcm_s16le_16" {
		t.Fatalf("pcm = %s", FileFormatPCMS16LE16)
	}
}

func TestSpeechToSpeechQuery(t *testing.T) {
	var opt speechToSpeechOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithoutSpeechToSpeechEnableLogging()(&opt)
	WithSpeechToSpeechOptimizeStreamingLatency(StreamingLatencyStrong)(&opt)
	WithSpeechToSpeechOutputFormat(MP3_44100_128)(&opt)

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

func TestSpeechToSpeechRequestBody(t *testing.T) {
	var opt speechToSpeechOption
	body, contentType, err := opt.body(strings.NewReader("source-audio"), "clip.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data; boundary=") {
		t.Fatalf("content type = %q", contentType)
	}

	form := parseMultipart(t, body, contentType)
	if got := formValue(t, form, "audio"); got != "source-audio" {
		t.Fatalf("audio = %q", got)
	}
	if form.File["audio"][0].Filename != "clip.mp3" {
		t.Fatalf("filename = %q", form.File["audio"][0].Filename)
	}

	WithSpeechToSpeechModelId("eleven_multilingual_sts_v2")(&opt)
	WithSpeechToSpeechVoiceSettings(VoiceSettings{Stability: 0.5, SimilarityBoost: 0.75, Speed: 1})(&opt)
	WithSpeechToSpeechSeed(42)(&opt)
	WithRemoveBackgroundNoise()(&opt)
	WithFileFormat(FileFormatPCMS16LE16)(&opt)

	body, contentType, err = opt.body(strings.NewReader("source-audio"), "clip.mp3")
	if err != nil {
		t.Fatal(err)
	}
	form = parseMultipart(t, body, contentType)
	if form.Value["model_id"][0] != "eleven_multilingual_sts_v2" {
		t.Fatalf("model_id = %v", form.Value["model_id"])
	}
	if form.Value["seed"][0] != "42" {
		t.Fatalf("seed = %v", form.Value["seed"])
	}
	if form.Value["remove_background_noise"][0] != "true" {
		t.Fatalf("remove_background_noise = %v", form.Value["remove_background_noise"])
	}
	if form.Value["file_format"][0] != "pcm_s16le_16" {
		t.Fatalf("file_format = %v", form.Value["file_format"])
	}

	var settings VoiceSettings
	if err := json.Unmarshal([]byte(form.Value["voice_settings"][0]), &settings); err != nil {
		t.Fatal(err)
	}
	if settings.Stability != 0.5 || settings.Speed != 1 {
		t.Fatalf("voice_settings = %+v", settings)
	}
}

func TestSpeechToSpeechValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.SpeechToSpeech(nil, "voice"); err == nil || err.Error() != "audio is nil" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToSpeech(strings.NewReader("audio"), ""); err == nil || err.Error() != "voiceId is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToSpeech(strings.NewReader(""), "voice"); err == nil || err.Error() != "audio is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToSpeechStream(nil, "voice"); err == nil || err.Error() != "audio is nil" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToSpeechStream(strings.NewReader("audio"), ""); err == nil || err.Error() != "voiceId is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestSpeechToSpeech(t *testing.T) {
	audio := []byte("ID3fake-mp3")
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/voice-changer/voice-1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("output_format") != "mp3_44100_128" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("content-type = %q", r.Header.Get("Content-Type"))
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("model_id") != "eleven_multilingual_sts_v2" {
			t.Fatalf("model_id = %q", r.FormValue("model_id"))
		}
		file, header, err := r.FormFile("audio")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		if header.Filename != "source.mp3" {
			t.Fatalf("filename = %q", header.Filename)
		}
		got, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "source-audio" {
			t.Fatalf("audio = %q", got)
		}
		w.Write(audio)
	})

	got, err := client.SpeechToSpeech(
		strings.NewReader("source-audio"),
		"voice-1",
		WithSpeechToSpeechModelId("eleven_multilingual_sts_v2"),
		WithSpeechToSpeechOutputFormat(MP3_44100_128),
		WithAudioFilename("source.mp3"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(audio) {
		t.Fatalf("audio = %q", got)
	}
}

func TestSpeechToSpeechError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad voice", http.StatusBadRequest)
	})
	_, err := client.SpeechToSpeech(strings.NewReader("audio"), "voice-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSpeechToSpeechStream(t *testing.T) {
	audio := []byte("streamed-mp3-bytes")
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/voice-changer/voice-1/stream" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write(audio)
	})

	stream, err := client.SpeechToSpeechStream(
		strings.NewReader("source-audio"),
		"voice-1",
		WithSpeechToSpeechOutputFormat(MP3_44100_128),
	)
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

func TestSpeechToSpeechStreamError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "wav not supported", http.StatusBadRequest)
	})
	_, err := client.SpeechToSpeechStream(strings.NewReader("audio"), "voice-1", WithSpeechToSpeechOutputFormat(WAV_44100))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAudioFilename(t *testing.T) {
	if got := audioFilename(strings.NewReader("x")); got != "audio" {
		t.Fatalf("reader = %q", got)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "take-1.wav")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if got := audioFilename(f); got != "take-1.wav" {
		t.Fatalf("file = %q", got)
	}
}

func TestLiveSpeechToSpeech(t *testing.T) {
	client := liveClient(t)
	source, err := client.TextToSpeech(
		"Hello from go eleven.",
		"CwhRBWXzGAHq8TQ4Fs17",
		WithModelId("eleven_multilingual_v2"),
		WithOutputFormat(MP3_44100_128),
	)
	if err != nil {
		t.Fatal(err)
	}

	audio, err := client.SpeechToSpeech(
		bytes.NewReader(source),
		"CwhRBWXzGAHq8TQ4Fs17",
		WithSpeechToSpeechModelId("eleven_multilingual_sts_v2"),
		WithSpeechToSpeechOutputFormat(MP3_44100_128),
		WithAudioFilename("source.mp3"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(audio) == 0 {
		t.Fatal("empty audio")
	}
	t.Logf("bytes=%d header=%q", len(audio), audio[:min(3, len(audio))])
}

func TestLiveSpeechToSpeechStream(t *testing.T) {
	client := liveClient(t)
	source, err := client.TextToSpeech(
		"Hello from go eleven.",
		"CwhRBWXzGAHq8TQ4Fs17",
		WithModelId("eleven_multilingual_v2"),
		WithOutputFormat(MP3_44100_128),
	)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := client.SpeechToSpeechStream(
		bytes.NewReader(source),
		"CwhRBWXzGAHq8TQ4Fs17",
		WithSpeechToSpeechModelId("eleven_multilingual_sts_v2"),
		WithSpeechToSpeechOutputFormat(MP3_44100_128),
		WithAudioFilename("source.mp3"),
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

func parseMultipart(t *testing.T, body []byte, contentType string) *multipart.Form {
	t.Helper()
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatal(err)
	}
	r := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { form.RemoveAll() })
	return form
}

func formValue(t *testing.T, form *multipart.Form, key string) string {
	t.Helper()
	files := form.File[key]
	if len(files) == 0 {
		t.Fatalf("missing file %q", key)
	}
	f, err := files[0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
