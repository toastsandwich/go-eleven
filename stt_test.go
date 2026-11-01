package eleven

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSpeechToTextEnums(t *testing.T) {
	if TimestampsGranularityWord.String() != "word" {
		t.Fatalf("granularity = %s", TimestampsGranularityWord)
	}
	if MultichannelOutputCombined.String() != "combined" {
		t.Fatalf("style = %s", MultichannelOutputCombined)
	}
	if EntityRedactionEnumeratedEntityType.String() != "enumerated_entity_type" {
		t.Fatalf("redaction = %s", EntityRedactionEnumeratedEntityType)
	}
	if SpeechToTextWordTypeAudioEvent.String() != "audio_event" {
		t.Fatalf("word type = %s", SpeechToTextWordTypeAudioEvent)
	}
	if ExportFormatSegmentedJSON.String() != "segmented_json" {
		t.Fatalf("export = %s", ExportFormatSegmentedJSON)
	}
}

func TestSpeechToTextQuery(t *testing.T) {
	var opt speechToTextOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithSpeechToTextToken("tok_1")(&opt)
	WithoutSpeechToTextEnableLogging()(&opt)

	q := opt.query()
	if q.Get("token") != "tok_1" {
		t.Fatalf("token = %q", q.Get("token"))
	}
	if q.Get("enable_logging") != "false" {
		t.Fatalf("enable_logging = %q", q.Get("enable_logging"))
	}
}

func TestSpeechToTextRequestBody(t *testing.T) {
	var opt speechToTextOption
	opt.ModelId = "scribe_v2"
	body, contentType, err := opt.body(strings.NewReader("source-audio"), "clip.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data; boundary=") {
		t.Fatalf("content type = %q", contentType)
	}

	form := parseMultipart(t, body, contentType)
	if got := formValue(t, form, "file"); got != "source-audio" {
		t.Fatalf("file = %q", got)
	}
	if form.File["file"][0].Filename != "clip.mp3" {
		t.Fatalf("filename = %q", form.File["file"][0].Filename)
	}
	if form.Value["model_id"][0] != "scribe_v2" {
		t.Fatalf("model_id = %v", form.Value["model_id"])
	}

	WithSpeechToTextLanguageCode("eng")(&opt)
	WithTranscriptEdit("make it concise")(&opt)
	WithTagAudioEvents()(&opt)
	WithNumSpeakers(2)(&opt)
	WithTimestampsGranularity(TimestampsGranularityCharacter)(&opt)
	WithDiarize()(&opt)
	WithDiarizationThreshold(0.22)(&opt)
	speakers := true
	WithAdditionalFormats(TranscriptExportOptions{Format: ExportFormatSRT, IncludeSpeakers: &speakers})(&opt)
	WithSpeechToTextFileFormat(FileFormatPCMS16LE16)(&opt)
	WithSpeechToTextWebhook()(&opt)
	WithSpeechToTextWebhookID("wh_1")(&opt)
	WithSpeechToTextTemperature(0.5)(&opt)
	WithSpeechToTextSeed(7)(&opt)
	WithUseMultiChannel()(&opt)
	WithMultichannelOutputStyle(MultichannelOutputCombined)(&opt)
	WithSpeechToTextWebhookMetadata(map[string]string{"job": "42"})(&opt)
	WithEntityDetection("all")(&opt)
	WithNoVerbatim()(&opt)
	WithUseSpeakerLibrary()(&opt)
	WithDetectSpeakerRoles()(&opt)
	WithEntityRedaction("pii", "phi")(&opt)
	WithEntityRedactionMode(EntityRedactionRedacted)(&opt)
	WithKeyterms("ElevenLabs", "Scribe")(&opt)

	body, contentType, err = opt.body(strings.NewReader("source-audio"), "clip.mp3")
	if err != nil {
		t.Fatal(err)
	}
	form = parseMultipart(t, body, contentType)
	if form.Value["language_code"][0] != "eng" {
		t.Fatalf("language_code = %v", form.Value["language_code"])
	}
	if form.Value["transcript_edit"][0] != "make it concise" {
		t.Fatalf("transcript_edit = %v", form.Value["transcript_edit"])
	}
	if form.Value["tag_audio_events"][0] != "true" {
		t.Fatalf("tag_audio_events = %v", form.Value["tag_audio_events"])
	}
	if form.Value["num_speakers"][0] != "2" {
		t.Fatalf("num_speakers = %v", form.Value["num_speakers"])
	}
	if form.Value["timestamps_granularity"][0] != "character" {
		t.Fatalf("timestamps_granularity = %v", form.Value["timestamps_granularity"])
	}
	if form.Value["diarize"][0] != "true" {
		t.Fatalf("diarize = %v", form.Value["diarize"])
	}
	if form.Value["diarization_threshold"][0] != "0.22" {
		t.Fatalf("diarization_threshold = %v", form.Value["diarization_threshold"])
	}
	if form.Value["file_format"][0] != "pcm_s16le_16" {
		t.Fatalf("file_format = %v", form.Value["file_format"])
	}
	if form.Value["webhook"][0] != "true" {
		t.Fatalf("webhook = %v", form.Value["webhook"])
	}
	if form.Value["webhook_id"][0] != "wh_1" {
		t.Fatalf("webhook_id = %v", form.Value["webhook_id"])
	}
	if form.Value["temperature"][0] != "0.5" {
		t.Fatalf("temperature = %v", form.Value["temperature"])
	}
	if form.Value["seed"][0] != "7" {
		t.Fatalf("seed = %v", form.Value["seed"])
	}
	if form.Value["use_multi_channel"][0] != "true" {
		t.Fatalf("use_multi_channel = %v", form.Value["use_multi_channel"])
	}
	if form.Value["multichannel_output_style"][0] != "combined" {
		t.Fatalf("multichannel_output_style = %v", form.Value["multichannel_output_style"])
	}
	if form.Value["entity_detection"][0] != "all" {
		t.Fatalf("entity_detection = %v", form.Value["entity_detection"])
	}
	if form.Value["no_verbatim"][0] != "true" {
		t.Fatalf("no_verbatim = %v", form.Value["no_verbatim"])
	}
	if form.Value["use_speaker_library"][0] != "true" {
		t.Fatalf("use_speaker_library = %v", form.Value["use_speaker_library"])
	}
	if form.Value["detect_speaker_roles"][0] != "true" {
		t.Fatalf("detect_speaker_roles = %v", form.Value["detect_speaker_roles"])
	}
	if form.Value["entity_redaction"][0] != `["pii","phi"]` {
		t.Fatalf("entity_redaction = %v", form.Value["entity_redaction"])
	}
	if form.Value["entity_redaction_mode"][0] != "redacted" {
		t.Fatalf("entity_redaction_mode = %v", form.Value["entity_redaction_mode"])
	}
	if got := form.Value["keyterms"]; len(got) != 2 || got[0] != "ElevenLabs" || got[1] != "Scribe" {
		t.Fatalf("keyterms = %v", got)
	}

	var formats []TranscriptExportOptions
	if err := json.Unmarshal([]byte(form.Value["additional_formats"][0]), &formats); err != nil {
		t.Fatal(err)
	}
	if len(formats) != 1 || formats[0].Format != ExportFormatSRT || formats[0].IncludeSpeakers == nil || !*formats[0].IncludeSpeakers {
		t.Fatalf("additional_formats = %+v", formats)
	}

	var meta map[string]string
	if err := json.Unmarshal([]byte(form.Value["webhook_metadata"][0]), &meta); err != nil {
		t.Fatal(err)
	}
	if meta["job"] != "42" {
		t.Fatalf("webhook_metadata = %v", meta)
	}
}

func TestSpeechToTextSourceURLBody(t *testing.T) {
	var opt speechToTextOption
	opt.ModelId = "scribe_v2"
	WithSourceURL("https://example.com/clip.mp3")(&opt)

	body, contentType, err := opt.body(nil, "audio")
	if err != nil {
		t.Fatal(err)
	}
	form := parseMultipart(t, body, contentType)
	if _, ok := form.File["file"]; ok {
		t.Fatal("expected no file part")
	}
	if form.Value["source_url"][0] != "https://example.com/clip.mp3" {
		t.Fatalf("source_url = %v", form.Value["source_url"])
	}
}

func TestSpeechToTextValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.SpeechToText(nil); err == nil || err.Error() != "audio is nil" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToText(strings.NewReader("")); err == nil || err.Error() != "audio is empty" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SpeechToText(strings.NewReader("audio"), WithSourceURL("https://example.com/a.mp3")); err == nil || err.Error() != "provide audio or source URL, not both" {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.GetTranscript(""); err == nil || err.Error() != "transcriptionID is empty" {
		t.Fatalf("err = %v", err)
	}
	if err := client.DeleteTranscript(""); err == nil || err.Error() != "transcriptionID is empty" {
		t.Fatalf("err = %v", err)
	}
}

func TestSpeechToText(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/speech-to-text" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("enable_logging") != "false" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("content-type = %q", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("model_id") != "scribe_v2" {
			t.Fatalf("model_id = %q", r.FormValue("model_id"))
		}
		if r.FormValue("language_code") != "eng" {
			t.Fatalf("language_code = %q", r.FormValue("language_code"))
		}
		if r.FormValue("diarize") != "true" {
			t.Fatalf("diarize = %q", r.FormValue("diarize"))
		}
		file, header, err := r.FormFile("file")
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
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"language_code":"eng","language_probability":0.99,"text":"Hello from go eleven.","words":[{"text":"Hello","type":"word","logprob":-0.1,"start":0,"end":0.4}],"transcription_id":"tr_1","audio_duration_secs":1.2}`))
	})

	got, err := client.SpeechToText(
		strings.NewReader("source-audio"),
		WithoutSpeechToTextEnableLogging(),
		WithSpeechToTextLanguageCode("eng"),
		WithDiarize(),
		WithSpeechToTextAudioFilename("source.mp3"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "Hello from go eleven." || got.LanguageCode != "eng" || got.TranscriptionID != "tr_1" {
		t.Fatalf("result = %+v", got)
	}
	if len(got.Words) != 1 || got.Words[0].Text != "Hello" || got.Words[0].Type != SpeechToTextWordTypeWord {
		t.Fatalf("words = %+v", got.Words)
	}
	if got.AudioDurationSecs == nil || *got.AudioDurationSecs != 1.2 {
		t.Fatalf("duration = %v", got.AudioDurationSecs)
	}
	if got.IsWebhook() || got.IsMultichannel() {
		t.Fatalf("shape webhook=%v multi=%v", got.IsWebhook(), got.IsMultichannel())
	}
}

func TestSpeechToTextDefaultModel(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("model_id") != "scribe_v2" {
			t.Fatalf("model_id = %q", r.FormValue("model_id"))
		}
		w.Write([]byte(`{"language_code":"eng","language_probability":1,"text":"ok","words":[]}`))
	})
	if _, err := client.SpeechToText(strings.NewReader("audio")); err != nil {
		t.Fatal(err)
	}
}

func TestSpeechToTextSourceURL(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("source_url") != "https://example.com/clip.mp3" {
			t.Fatalf("source_url = %q", r.FormValue("source_url"))
		}
		if r.FormValue("model_id") != "scribe_v2_medical" {
			t.Fatalf("model_id = %q", r.FormValue("model_id"))
		}
		if _, _, err := r.FormFile("file"); err == nil {
			t.Fatal("expected no file")
		}
		w.Write([]byte(`{"language_code":"eng","language_probability":1,"text":"from url","words":[]}`))
	})

	got, err := client.SpeechToText(nil, WithSourceURL("https://example.com/clip.mp3"), WithSpeechToTextModelId("scribe_v2_medical"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "from url" {
		t.Fatalf("text = %q", got.Text)
	}
}

func TestSpeechToTextWebhook(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message":"accepted","request_id":"req_1","transcription_id":"tr_async"}`))
	})

	got, err := client.SpeechToText(strings.NewReader("audio"), WithSpeechToTextWebhook())
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsWebhook() || got.RequestID != "req_1" || got.TranscriptionID != "tr_async" {
		t.Fatalf("result = %+v", got)
	}
}

func TestSpeechToTextMultichannel(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"transcription_id":"tr_m","audio_duration_secs":3,"transcripts":[{"text":"left","language_code":"eng","language_probability":1,"words":[],"channel_index":0},{"text":"right","language_code":"eng","language_probability":1,"words":[],"channel_index":1}]}`))
	})

	got, err := client.SpeechToText(strings.NewReader("audio"), WithUseMultiChannel())
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsMultichannel() || got.TranscriptionID != "tr_m" {
		t.Fatalf("result = %+v", got)
	}
	var texts []string
	for transcript := range got.All() {
		texts = append(texts, transcript.Text)
	}
	if len(texts) != 2 || texts[0] != "left" || texts[1] != "right" {
		t.Fatalf("texts = %v", texts)
	}
}

func TestSpeechToTextError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad audio", http.StatusBadRequest)
	})
	_, err := client.SpeechToText(strings.NewReader("audio"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetTranscript(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/speech-to-text/transcripts/tr_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		w.Write([]byte(`{"language_code":"eng","language_probability":0.9,"text":"stored","words":[],"transcription_id":"tr_1"}`))
	})

	got, err := client.GetTranscript("tr_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "stored" || got.TranscriptionID != "tr_1" {
		t.Fatalf("result = %+v", got)
	}
}

func TestGetTranscriptError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if _, err := client.GetTranscript("tr_missing"); err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteTranscript(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/speech-to-text/transcripts/tr_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.DeleteTranscript("tr_1"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTranscriptError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if err := client.DeleteTranscript("tr_missing"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSpeechToTextResultAll(t *testing.T) {
	single := SpeechToTextResult{SpeechToTextTranscript: SpeechToTextTranscript{Text: "one"}}
	var texts []string
	for transcript := range single.All() {
		texts = append(texts, transcript.Text)
	}
	if len(texts) != 1 || texts[0] != "one" {
		t.Fatalf("single = %v", texts)
	}

	multi := SpeechToTextResult{
		Transcripts: SpeechToTextTranscripts{
			{Text: "a", LanguageCode: "eng"},
			{Text: "b", LanguageCode: "spa"},
		},
	}
	var filtered []string
	for transcript := range multi.AllWith(func(tr *SpeechToTextTranscript) bool {
		return tr.LanguageCode == "spa"
	}) {
		filtered = append(filtered, transcript.Text)
	}
	if len(filtered) != 1 || filtered[0] != "b" {
		t.Fatalf("filtered = %v", filtered)
	}

	webhook := SpeechToTextResult{Message: "accepted", RequestID: "req_1"}
	for range webhook.All() {
		t.Fatal("webhook All yielded a transcript")
	}
}

func TestSpeechToTextTranscriptionDataResult(t *testing.T) {
	data := SpeechToTextTranscriptionData{
		RequestID:     "req",
		Transcription: json.RawMessage(`{"text":"hello","language_code":"eng","language_probability":1,"words":[]}`),
	}
	result, err := data.Result()
	if err != nil {
		t.Fatal(err)
	}
	if result.Text != "hello" {
		t.Fatalf("text = %q", result.Text)
	}

	if _, err := (SpeechToTextTranscriptionData{}).Result(); err == nil || err.Error() != "webhook event has no transcription" {
		t.Fatalf("err = %v", err)
	}
}

func TestEditedTranscriptOK(t *testing.T) {
	if !(EditedTranscript{Kind: "transcript"}).OK() {
		t.Fatal("expected ok")
	}
	if (EditedTranscript{Kind: "error"}).OK() {
		t.Fatal("expected not ok")
	}
}

func TestWithoutTagAudioEvents(t *testing.T) {
	var opt speechToTextOption
	WithoutTagAudioEvents()(&opt)
	if opt.TagAudioEvents == nil || *opt.TagAudioEvents {
		t.Fatal("expected false")
	}
}

func TestLiveSpeechToText(t *testing.T) {
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

	result, err := client.SpeechToText(
		bytes.NewReader(source),
		WithSpeechToTextModelId("scribe_v2"),
		WithSpeechToTextAudioFilename("source.mp3"),
		WithDiarize(),
		WithTagAudioEvents(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Text == "" {
		t.Fatal("empty transcript")
	}
	t.Logf("text=%q language=%s id=%s", result.Text, result.LanguageCode, result.TranscriptionID)

	if result.TranscriptionID != "" {
		stored, err := client.GetTranscript(result.TranscriptionID)
		if err != nil {
			t.Fatal(err)
		}
		if stored.Text == "" {
			t.Fatal("empty stored transcript")
		}
		if err := client.DeleteTranscript(result.TranscriptionID); err != nil {
			t.Fatal(err)
		}
	}
}
