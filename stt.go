package eleven

import (
	"encoding/json"
	"fmt"
	"iter"
)

type TimestampsGranularity string

const (
	TimestampsGranularityNone      TimestampsGranularity = "none"
	TimestampsGranularityWord      TimestampsGranularity = "word"
	TimestampsGranularityCharacter TimestampsGranularity = "character"
)

func (g TimestampsGranularity) String() string { return string(g) }

type MultichannelOutputStyle string

const (
	MultichannelOutputSeparate MultichannelOutputStyle = "separate"
	MultichannelOutputCombined MultichannelOutputStyle = "combined"
)

func (s MultichannelOutputStyle) String() string { return string(s) }

type EntityRedactionMode string

const (
	EntityRedactionRedacted             EntityRedactionMode = "redacted"
	EntityRedactionEntityType           EntityRedactionMode = "entity_type"
	EntityRedactionEnumeratedEntityType EntityRedactionMode = "enumerated_entity_type"
)

func (m EntityRedactionMode) String() string { return string(m) }

type SpeechToTextWordType string

const (
	SpeechToTextWordTypeWord       SpeechToTextWordType = "word"
	SpeechToTextWordTypeSpacing    SpeechToTextWordType = "spacing"
	SpeechToTextWordTypeAudioEvent SpeechToTextWordType = "audio_event"
)

func (t SpeechToTextWordType) String() string { return string(t) }

type ExportFormat string

const (
	ExportFormatDocx          ExportFormat = "docx"
	ExportFormatHTML          ExportFormat = "html"
	ExportFormatPDF           ExportFormat = "pdf"
	ExportFormatSegmentedJSON ExportFormat = "segmented_json"
	ExportFormatSRT           ExportFormat = "srt"
	ExportFormatTXT           ExportFormat = "txt"
)

func (f ExportFormat) String() string { return string(f) }

type TranscriptExportOptions struct {
	Format                      ExportFormat `json:"format"`
	IncludeSpeakers             *bool        `json:"include_speakers,omitempty"`
	IncludeTimestamps           *bool        `json:"include_timestamps,omitempty"`
	SegmentOnSilenceLongerThanS *float64     `json:"segment_on_silence_longer_than_s,omitempty"`
	MaxSegmentDurationS         *float64     `json:"max_segment_duration_s,omitempty"`
	MaxSegmentChars             *int         `json:"max_segment_chars,omitempty"`
	MaxCharactersPerLine        *int         `json:"max_characters_per_line,omitempty"`
}

type SpeechToTextCharacter struct {
	Text  string   `json:"text"`
	Start *float64 `json:"start"`
	End   *float64 `json:"end"`
}

type SpeechToTextWord struct {
	Text         string                  `json:"text"`
	Start        *float64                `json:"start"`
	End          *float64                `json:"end"`
	Type         SpeechToTextWordType    `json:"type"`
	SpeakerID    string                  `json:"speaker_id,omitempty"`
	Logprob      float64                 `json:"logprob"`
	Characters   []SpeechToTextCharacter `json:"characters,omitempty"`
	ChannelIndex *int                    `json:"channel_index,omitempty"`
}

type DetectedEntity struct {
	Text       string `json:"text"`
	EntityType string `json:"entity_type"`
	StartChar  int    `json:"start_char"`
	EndChar    int    `json:"end_char"`
}

type AdditionalFormat struct {
	RequestedFormat string `json:"requested_format"`
	FileExtension   string `json:"file_extension"`
	ContentType     string `json:"content_type"`
	IsBase64Encoded bool   `json:"is_base64_encoded"`
	Content         string `json:"content"`
}

type EditedTranscript struct {
	Kind        string `json:"kind"`
	ErrorType   string `json:"error_type,omitempty"`
	Message     string `json:"message,omitempty"`
	MessageType string `json:"message_type,omitempty"`
	Text        string `json:"text,omitempty"`
	EditedText  string `json:"edited_text,omitempty"`
}

func (e EditedTranscript) OK() bool {
	return e.Kind == "transcript"
}

type SpeechToTextTranscript struct {
	LanguageCode        string             `json:"language_code,omitempty"`
	LanguageProbability float64            `json:"language_probability,omitempty"`
	Text                string             `json:"text,omitempty"`
	Words               []SpeechToTextWord `json:"words,omitempty"`
	ChannelIndex        *int               `json:"channel_index,omitempty"`
	AdditionalFormats   []AdditionalFormat `json:"additional_formats,omitempty"`
	TranscriptionID     string             `json:"transcription_id,omitempty"`
	Entities            []DetectedEntity   `json:"entities,omitempty"`
	AudioDurationSecs   *float64           `json:"audio_duration_secs,omitempty"`
	EditedTranscript    *EditedTranscript  `json:"edited_transcript,omitempty"`
}

type SpeechToTextTranscripts []SpeechToTextTranscript

type SpeechToTextFilterFunc func(*SpeechToTextTranscript) bool

func (t SpeechToTextTranscripts) All() iter.Seq[SpeechToTextTranscript] {
	return func(yield func(SpeechToTextTranscript) bool) {
		for _, transcript := range t {
			if !yield(transcript) {
				return
			}
		}
	}
}

func (t SpeechToTextTranscripts) AllWith(filters ...SpeechToTextFilterFunc) iter.Seq[SpeechToTextTranscript] {
	return func(yield func(SpeechToTextTranscript) bool) {
		for _, transcript := range t {
			if !matchAllSpeechToText(&transcript, filters) {
				continue
			}
			if !yield(transcript) {
				return
			}
		}
	}
}

func matchAllSpeechToText(transcript *SpeechToTextTranscript, filters []SpeechToTextFilterFunc) bool {
	for _, filter := range filters {
		if !filter(transcript) {
			return false
		}
	}
	return true
}

type SpeechToTextResult struct {
	SpeechToTextTranscript
	Transcripts SpeechToTextTranscripts `json:"transcripts,omitempty"`
	Message     string                  `json:"message,omitempty"`
	RequestID   string                  `json:"request_id,omitempty"`
}

func (r SpeechToTextResult) IsWebhook() bool {
	return r.RequestID != "" || r.Message != ""
}

func (r SpeechToTextResult) IsMultichannel() bool {
	return len(r.Transcripts) > 0
}

func (r SpeechToTextResult) All() iter.Seq[SpeechToTextTranscript] {
	if r.IsMultichannel() {
		return r.Transcripts.All()
	}
	if r.IsWebhook() {
		return SpeechToTextTranscripts{}.All()
	}
	return SpeechToTextTranscripts{r.SpeechToTextTranscript}.All()
}

func (r SpeechToTextResult) AllWith(filters ...SpeechToTextFilterFunc) iter.Seq[SpeechToTextTranscript] {
	if r.IsMultichannel() {
		return r.Transcripts.AllWith(filters...)
	}
	if r.IsWebhook() {
		return SpeechToTextTranscripts{}.AllWith(filters...)
	}
	return SpeechToTextTranscripts{r.SpeechToTextTranscript}.AllWith(filters...)
}

func (d SpeechToTextTranscriptionData) Result() (SpeechToTextResult, error) {
	if len(d.Transcription) == 0 {
		return SpeechToTextResult{}, fmt.Errorf("webhook event has no transcription")
	}
	var result SpeechToTextResult
	if err := json.Unmarshal(d.Transcription, &result); err != nil {
		return SpeechToTextResult{}, fmt.Errorf("failed to marshal transcription: %v", err)
	}
	return result, nil
}
