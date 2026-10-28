package eleven

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// ============================================================
// Text To Speech Options
// ============================================================

type OutputFormat byte

const (
	ALAW8000 OutputFormat = iota
	MP3_22050_32
	MP3_24000_48
	MP3_44100_32
	MP3_44100_64
	MP3_44100_96
	MP3_44100_128
	MP3_44100_192
	OPUS_48000_32
	OPUS_48000_64
	OPUS_48000_96
	OPUS_48000_128
	OPUS_48000_192
	PCM_8000
	PCM_16000
	PCM_22050
	PCM_24000
	PCM_32000
	PCM_44100
	PCM_48000
	ULAW8000
	WAV_8000
	WAV_16000
	WAV_22050
	WAV_24000
	WAV_32000
	WAV_44100
	WAV_48000
)

func (o OutputFormat) String() string {
	switch o {
	case ALAW8000:
		return "alaw_8000"
	case MP3_22050_32:
		return "mp3_22050_32"
	case MP3_24000_48:
		return "mp3_24000_48"
	case MP3_44100_32:
		return "mp3_44100_32"
	case MP3_44100_64:
		return "mp3_44100_64"
	case MP3_44100_96:
		return "mp3_44100_96"
	case MP3_44100_128:
		return "mp3_44100_128"
	case MP3_44100_192:
		return "mp3_44100_192"
	case OPUS_48000_32:
		return "opus_48000_32"
	case OPUS_48000_64:
		return "opus_48000_64"
	case OPUS_48000_96:
		return "opus_48000_96"
	case OPUS_48000_128:
		return "opus_48000_128"
	case OPUS_48000_192:
		return "opus_48000_192"
	case PCM_8000:
		return "pcm_8000"
	case PCM_16000:
		return "pcm_16000"
	case PCM_22050:
		return "pcm_22050"
	case PCM_24000:
		return "pcm_24000"
	case PCM_32000:
		return "pcm_32000"
	case PCM_44100:
		return "pcm_44100"
	case PCM_48000:
		return "pcm_48000"
	case ULAW8000:
		return "ulaw_8000"
	case WAV_8000:
		return "wav_8000"
	case WAV_16000:
		return "wav_16000"
	case WAV_22050:
		return "wav_22050"
	case WAV_24000:
		return "wav_24000"
	case WAV_32000:
		return "wav_32000"
	case WAV_44100:
		return "wav_44100"
	case WAV_48000:
		return "wav_48000"
	default:
		return "unknown"
	}
}

type StreamingLatency byte

const (
	StreamingLatencyDefault StreamingLatency = iota
	StreamingLatencyNormal
	StreamingLatencyStrong
	StreamingLatencyMax
	StreamingLatencyMaxNoNormalizer
)

func (l StreamingLatency) String() string {
	return strconv.Itoa(int(l))
}

type TextNormalization string

const (
	TextNormalizationAuto TextNormalization = "auto"
	TextNormalizationOn   TextNormalization = "on"
	TextNormalizationOff  TextNormalization = "off"
)

func (t TextNormalization) String() string { return string(t) }

type PronunciationDictionaryLocator struct {
	PronunciationDictionaryId string `json:"pronunciation_dictionary_id"`
	VersionId                 string `json:"version_id,omitempty"`
}

type AudioEffects struct {
	FilterPresetId    string  `json:"filter_preset_id,omitempty"`
	Distance          float64 `json:"distance,omitempty"`
	EnvironmentId     string  `json:"environment_id,omitempty"`
	BackgroundNoiseId string  `json:"background_noise_id,omitempty"`
	SendLevel         float64 `json:"send_level,omitempty"`
	Pan               float64 `json:"pan,omitempty"`
}

type textToSpeechOption struct {
	// queryParamters
	EnableLogging            *bool
	OptimizeStreamingLatency *int
	OutputFormat             string
	// body
	ModelId                         string
	LanguageCode                    string
	VoiceSettings                   *VoiceSettings
	PronunciationDictionaryLocators []PronunciationDictionaryLocator
	Seed                            *int
	PreviousText                    string
	NextText                        string
	PreviousRequestIds              []string
	NextRequestIds                  []string
	UsePvcAsIvc                     *bool
	ApplyTextNormalization          string
	ApplyLanguageTextNormalization  bool
	AudioEffects                    *AudioEffects
}

type TextToSpeechOption func(*textToSpeechOption)

// When enable_logging is set to false zero retention mode will be used for the request. This will mean history features are unavailable for this request, including request stitching. Zero retention mode may only be used by enterprise customers.
func WithoutEnableLogging() TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		enabled := false
		ttso.EnableLogging = &enabled
	}
}

// You can turn on latency optimizations at some cost of quality. The best possible final latency varies by model. Possible values: 0 - default mode (no latency optimizations) 1 - normal latency optimizations (about 50% of possible latency improvement of option 3) 2 - strong latency optimizations (about 75% of possible latency improvement of option 3) 3 - max latency optimizations 4 - max latency optimizations, but also with text normalizer turned off for even more latency savings (best latency, but can mispronounce eg numbers and dates).
func WithOptimizeStreamingLatency(latency StreamingLatency) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		v := int(latency)
		ttso.OptimizeStreamingLatency = &v
	}
}

// Output format of the generated audio. Formatted as codec_sample_rate_bitrate. So an mp3 with 22.05kHz sample rate at 32kbs is represented as mp3_22050_32. MP3 with 192kbps bitrate requires you to be subscribed to Creator tier or above. PCM and WAV formats with 44.1kHz sample rate requires you to be subscribed to Pro tier or above. Note that the μ-law format (sometimes written mu-law, often approximated as u-law) is commonly used for Twilio audio inputs.
func WithOutputFormat(format OutputFormat) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.OutputFormat = format.String()
	}
}

// Identifier of the model that will be used, you can query them using GET /v1/models. The model needs to have support for text to speech, you can check this using the can_do_text_to_speech property. Defaults to eleven_multilingual_v2.
func WithModelId(modelId string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.ModelId = modelId
	}
}

// Language code (ISO 639-1) used to enforce a language for the model and text normalization. If the model does not support the provided language code, it will be ignored. This parameter is not supported for multilingual_v2 models.
func WithLanguageCode(languageCode string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.LanguageCode = languageCode
	}
}

// Voice settings overriding stored settings for the given voice. They are applied only on the given request.
func WithVoiceSettings(voiceSettings VoiceSettings) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.VoiceSettings = &voiceSettings
	}
}

// A list of pronunciation dictionary locators (id, version_id) to be applied to the text. They will be applied in order. You may have up to 3 locators per request.
func WithPronunciationDictionaryLocators(locators ...PronunciationDictionaryLocator) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.PronunciationDictionaryLocators = locators
	}
}

// If specified, our system will make a best effort to sample deterministically, such that repeated requests with the same seed and parameters should return the same result. Determinism is not guaranteed. Must be integer between 0 and 4294967295.
func WithSeed(seed int) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.Seed = &seed
	}
}

// The text that came before the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation.
func WithPreviousText(text string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.PreviousText = text
	}
}

// The text that comes after the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation.
func WithNextText(text string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.NextText = text
	}
}

// A list of request_id of the samples that were generated before this generation. Can be used to improve the speech's continuity when splitting up a large task into multiple requests. The results will be best when the same model is used across the generations. In case both previous_text and previous_request_ids is send, previous_text will be ignored. A maximum of 3 request_ids can be send.
func WithPreviousRequestIds(ids ...string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.PreviousRequestIds = ids
	}
}

// A list of request_id of the samples that come after this generation. next_request_ids is especially useful for maintaining the speech's continuity when regenerating a sample that has had some audio quality issues. For example, if you have generated 3 speech clips, and you want to improve clip 2, passing the request id of clip 3 as a next_request_id (and that of clip 1 as a previous_request_id) will help maintain natural flow in the combined speech. The results will be best when the same model is used across the generations. In case both next_text and next_request_ids is send, next_text will be ignored. A maximum of 3 request_ids can be send.
func WithNextRequestIds(ids ...string) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.NextRequestIds = ids
	}
}

// Whether to use the IVC version of a professional voice. This may improve expressiveness and reduce latency.
func WithUsePvcAsIvc(use bool) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.UsePvcAsIvc = &use
	}
}

// This parameter controls text normalization with three modes: 'auto', 'on', and 'off'. When set to 'auto', the system will automatically decide whether to apply text normalization (e.g., spelling out numbers). With 'on', text normalization will always be applied, while with 'off', it will be skipped.
func WithApplyTextNormalization(mode TextNormalization) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.ApplyTextNormalization = mode.String()
	}
}

// This parameter controls language text normalization. This helps with proper pronunciation of text in some supported languages. WARNING: This parameter can heavily increase the latency of the request. Currently only supported for Japanese.
func WithApplyLanguageTextNormalization() TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.ApplyLanguageTextNormalization = true
	}
}

// Audio effects applied to the generated audio. Unknown IDs return 404 with param naming the field. Reverb (environment_id with send_level above 0) and a non-zero pan produce stereo audio and are rejected with 400 for ulaw_*, alaw_*, and pcm_* output formats. Effects don't change the cost of the request.
func WithAudioEffects(effects AudioEffects) TextToSpeechOption {
	return func(ttso *textToSpeechOption) {
		ttso.AudioEffects = &effects
	}
}

func (o *textToSpeechOption) query() url.Values {
	q := url.Values{}
	if o.EnableLogging != nil {
		q.Set("enable_logging", strconv.FormatBool(*o.EnableLogging))
	}
	if o.OptimizeStreamingLatency != nil {
		q.Set("optimize_streaming_latency", strconv.Itoa(*o.OptimizeStreamingLatency))
	}
	if o.OutputFormat != "" {
		q.Set("output_format", o.OutputFormat)
	}
	return q
}

type textToSpeechRequest struct {
	Text                            string                           `json:"text"`
	ModelId                         string                           `json:"model_id,omitempty"`
	LanguageCode                    string                           `json:"language_code,omitempty"`
	VoiceSettings                   *VoiceSettings                   `json:"voice_settings,omitempty"`
	PronunciationDictionaryLocators []PronunciationDictionaryLocator `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            *int                             `json:"seed,omitempty"`
	PreviousText                    string                           `json:"previous_text,omitempty"`
	NextText                        string                           `json:"next_text,omitempty"`
	PreviousRequestIds              []string                         `json:"previous_request_ids,omitempty"`
	NextRequestIds                  []string                         `json:"next_request_ids,omitempty"`
	UsePvcAsIvc                     *bool                            `json:"use_pvc_as_ivc,omitempty"`
	ApplyTextNormalization          string                           `json:"apply_text_normalization,omitempty"`
	ApplyLanguageTextNormalization  bool                             `json:"apply_language_text_normalization,omitempty"`
	AudioEffects                    *AudioEffects                    `json:"audio_effects,omitempty"`
}

func (o *textToSpeechOption) body(text string) textToSpeechRequest {
	return textToSpeechRequest{
		Text:                            text,
		ModelId:                         o.ModelId,
		LanguageCode:                    o.LanguageCode,
		VoiceSettings:                   o.VoiceSettings,
		PronunciationDictionaryLocators: o.PronunciationDictionaryLocators,
		Seed:                            o.Seed,
		PreviousText:                    o.PreviousText,
		NextText:                        o.NextText,
		PreviousRequestIds:              o.PreviousRequestIds,
		NextRequestIds:                  o.NextRequestIds,
		UsePvcAsIvc:                     o.UsePvcAsIvc,
		ApplyTextNormalization:          o.ApplyTextNormalization,
		ApplyLanguageTextNormalization:  o.ApplyLanguageTextNormalization,
		AudioEffects:                    o.AudioEffects,
	}
}

// ============================================================
// Search Voices Options
// ============================================================

type searchVoicesOption struct {
	NextPageToken        string
	PageSize             int
	Search               string
	Sort                 string
	SortDirection        string
	VoiceType            string
	Category             string
	FineTuningState      string
	CollectionId         string
	Gender               string
	Age                  string
	Language             []string
	Accent               string
	UseCases             []string
	MinNoticePeriodDays  int
	IncludeCustomRates   *bool
	IncludeLiveModerated *bool
	HighQuality          *bool
	IncludeTotalCount    *bool
	VoiceIds             []string
}

type SearchVoicesOption func(*searchVoicesOption)

// The next page token to use for pagination. Returned from the previous request. Use this in combination with the has_more flag for reliable pagination.
func WithNextPageToken(token string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.NextPageToken = token
	}
}

// How many voices to return at maximum. Can not exceed 100, defaults to 10. Page 0 may include more voices due to default voices being included.
func WithPageSize(size int) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.PageSize = size
	}
}

// Search term to filter voices by. Searches in name, description, labels, category.
func WithSearch(search string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Search = search
	}
}

// Which field to sort by, one of 'created_at_unix' or 'name'. 'created_at_unix' may not be available for older voices.
func WithSort(sort VoiceSort) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Sort = sort.String()
	}
}

// Which direction to sort the voices in. 'asc' or 'desc'.
func WithSortDirection(direction SortDirection) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.SortDirection = direction.String()
	}
}

// Type of the voice to filter by. One of 'personal', 'community', 'default', 'workspace', 'non-default', 'non-community', 'saved'.
func WithVoiceType(voiceType VoiceType) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.VoiceType = voiceType.String()
	}
}

// Category of the voice to filter by. One of 'premade', 'cloned', 'generated', 'professional'.
func WithCategory(category VoiceCategory) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Category = category.String()
	}
}

// State of the voice's fine tuning to filter by. Applicable only to professional voices clones. One of 'draft', 'not_verified', 'not_started', 'queued', 'fine_tuning', 'fine_tuned', 'failed', 'delayed'.
func WithFineTuningState(state FineTuningState) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.FineTuningState = state.String()
	}
}

// Collection ID to filter voices by.
func WithCollectionId(collectionId string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.CollectionId = collectionId
	}
}

// Gender used for filtering, based on the voice's 'gender' label.
func WithGender(gender string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Gender = gender
	}
}

// Age used for filtering, based on the voice's 'age' label.
func WithAge(age string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Age = age
	}
}

// Languages used for filtering, based on the voice's 'language' label. Voices matching any of the given languages are returned.
func WithLanguages(languages ...string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Language = languages
	}
}

// Accent used for filtering, based on the voice's 'accent' label.
func WithAccent(accent string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.Accent = accent
	}
}

// Use cases used for filtering, based on the voice's 'use_case' label. Voices matching any of the given use cases are returned.
func WithUseCases(useCases ...string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.UseCases = useCases
	}
}

// Filter to voices whose sharing notice period is at least the given number of days.
func WithMinNoticePeriodDays(days int) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.MinNoticePeriodDays = days
	}
}

// Whether to include voices that have a custom sharing rate. Defaults to including them.
func WithIncludeCustomRates(include bool) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.IncludeCustomRates = &include
	}
}

// Whether to include voices that have live moderation enabled. Defaults to including them.
func WithIncludeLiveModerated(include bool) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.IncludeLiveModerated = &include
	}
}

// When true, only return studio-quality voices (those whose category is 'high_quality').
func WithHighQuality(highQuality bool) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.HighQuality = &highQuality
	}
}

// Whether to include the total count of voices found in the response. NOTE: The total_count value is a live snapshot and may change between requests. For pagination, rely on the has_more flag instead.
func WithIncludeTotalCount(include bool) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.IncludeTotalCount = &include
	}
}

// Voice IDs to lookup by. Maximum 100 voice IDs.
func WithVoiceIds(voiceIds ...string) SearchVoicesOption {
	return func(o *searchVoicesOption) {
		o.VoiceIds = voiceIds
	}
}

func (o *searchVoicesOption) query() url.Values {
	q := url.Values{}
	if o.NextPageToken != "" {
		q.Set("next_page_token", o.NextPageToken)
	}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.Search != "" {
		q.Set("search", o.Search)
	}
	if o.Sort != "" {
		q.Set("sort", o.Sort)
	}
	if o.SortDirection != "" {
		q.Set("sort_direction", o.SortDirection)
	}
	if o.VoiceType != "" {
		q.Set("voice_type", o.VoiceType)
	}
	if o.Category != "" {
		q.Set("category", o.Category)
	}
	if o.FineTuningState != "" {
		q.Set("fine_tuning_state", o.FineTuningState)
	}
	if o.CollectionId != "" {
		q.Set("collection_id", o.CollectionId)
	}
	if o.Gender != "" {
		q.Set("gender", o.Gender)
	}
	if o.Age != "" {
		q.Set("age", o.Age)
	}
	if len(o.Language) > 0 {
		q.Set("language", marshalQueryList(o.Language))
	}
	if o.Accent != "" {
		q.Set("accent", o.Accent)
	}
	if len(o.UseCases) > 0 {
		q.Set("use_cases", marshalQueryList(o.UseCases))
	}
	if o.MinNoticePeriodDays > 0 {
		q.Set("min_notice_period_days", strconv.Itoa(o.MinNoticePeriodDays))
	}
	if o.IncludeCustomRates != nil {
		q.Set("include_custom_rates", strconv.FormatBool(*o.IncludeCustomRates))
	}
	if o.IncludeLiveModerated != nil {
		q.Set("include_live_moderated", strconv.FormatBool(*o.IncludeLiveModerated))
	}
	if o.HighQuality != nil {
		q.Set("high_quality", strconv.FormatBool(*o.HighQuality))
	}
	if o.IncludeTotalCount != nil {
		q.Set("include_total_count", strconv.FormatBool(*o.IncludeTotalCount))
	}
	if len(o.VoiceIds) > 0 {
		q.Set("voice_ids", marshalQueryList(o.VoiceIds))
	}
	return q
}

func marshalQueryList(values []string) string {
	b, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(b)
}
