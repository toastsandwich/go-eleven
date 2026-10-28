package eleven

import "iter"

type VoiceCategory string

const (
	VoiceCategoryGenerated    VoiceCategory = "generated"
	VoiceCategoryCloned       VoiceCategory = "cloned"
	VoiceCategoryPremade      VoiceCategory = "premade"
	VoiceCategoryProfessional VoiceCategory = "professional"
	VoiceCategoryFamous       VoiceCategory = "famous"
	VoiceCategoryHighQuality  VoiceCategory = "high_quality"
)

func (c VoiceCategory) String() string { return string(c) }

type VoiceType string

const (
	VoiceTypePersonal     VoiceType = "personal"
	VoiceTypeCommunity    VoiceType = "community"
	VoiceTypeDefault      VoiceType = "default"
	VoiceTypeWorkspace    VoiceType = "workspace"
	VoiceTypeNonDefault   VoiceType = "non-default"
	VoiceTypeNonCommunity VoiceType = "non-community"
	VoiceTypeSaved        VoiceType = "saved"
)

func (t VoiceType) String() string { return string(t) }

type FineTuningState string

const (
	FineTuningStateDraft       FineTuningState = "draft"
	FineTuningStateNotVerified FineTuningState = "not_verified"
	FineTuningStateNotStarted  FineTuningState = "not_started"
	FineTuningStateQueued      FineTuningState = "queued"
	FineTuningStateFineTuning  FineTuningState = "fine_tuning"
	FineTuningStateFineTuned   FineTuningState = "fine_tuned"
	FineTuningStateFailed      FineTuningState = "failed"
	FineTuningStateDelayed     FineTuningState = "delayed"
)

func (s FineTuningState) String() string { return string(s) }

type VoiceSort string

const (
	VoiceSortCreatedAtUnix VoiceSort = "created_at_unix"
	VoiceSortName          VoiceSort = "name"
)

func (s VoiceSort) String() string { return string(s) }

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

func (d SortDirection) String() string { return string(d) }

type SafetyControl string

const (
	SafetyControlNone              SafetyControl = "NONE"
	SafetyControlBan               SafetyControl = "BAN"
	SafetyControlCaptcha           SafetyControl = "CAPTCHA"
	SafetyControlEnterpriseBan     SafetyControl = "ENTERPRISE_BAN"
	SafetyControlEnterpriseCaptcha SafetyControl = "ENTERPRISE_CAPTCHA"
)

func (s SafetyControl) String() string { return string(s) }

type RecordingQuality string

const (
	RecordingQualityStudio RecordingQuality = "studio"
	RecordingQualityGood   RecordingQuality = "good"
	RecordingQualityOK     RecordingQuality = "ok"
	RecordingQualityPoor   RecordingQuality = "poor"
	RecordingQualityBad    RecordingQuality = "bad"
)

func (q RecordingQuality) String() string { return string(q) }

type LabellingStatus string

const (
	LabellingStatusInReview       LabellingStatus = "in_review"
	LabellingStatusReviewComplete LabellingStatus = "review_complete"
)

func (s LabellingStatus) String() string { return string(s) }

type VoiceSharingStatus string

const (
	VoiceSharingEnabled        VoiceSharingStatus = "enabled"
	VoiceSharingDisabled       VoiceSharingStatus = "disabled"
	VoiceSharingCopied         VoiceSharingStatus = "copied"
	VoiceSharingCopiedDisabled VoiceSharingStatus = "copied_disabled"
)

func (s VoiceSharingStatus) String() string { return string(s) }

type ReviewStatus string

const (
	ReviewStatusNotRequested          ReviewStatus = "not_requested"
	ReviewStatusPending               ReviewStatus = "pending"
	ReviewStatusDeclined              ReviewStatus = "declined"
	ReviewStatusDeclinedResubmittable ReviewStatus = "declined_resubmittable"
	ReviewStatusAllowed               ReviewStatus = "allowed"
	ReviewStatusAllowedWithChanges    ReviewStatus = "allowed_with_changes"
)

func (s ReviewStatus) String() string { return string(s) }

type SpeakerSeparationStatus string

const (
	SpeakerSeparationNotStarted SpeakerSeparationStatus = "not_started"
	SpeakerSeparationPending    SpeakerSeparationStatus = "pending"
	SpeakerSeparationCompleted  SpeakerSeparationStatus = "completed"
	SpeakerSeparationFailed     SpeakerSeparationStatus = "failed"
)

func (s SpeakerSeparationStatus) String() string { return string(s) }

type ReaderResourceType string

const (
	ReaderResourceRead       ReaderResourceType = "read"
	ReaderResourceCollection ReaderResourceType = "collection"
)

func (t ReaderResourceType) String() string { return string(t) }

type Utterance struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type Speaker struct {
	SpeakerId    string      `json:"speaker_id"`
	DurationSecs float64     `json:"duration_secs"`
	Utterances   []Utterance `json:"utterances"`
}

type SpeakerSeparation struct {
	VoiceId            string                  `json:"voice_id"`
	SampleId           string                  `json:"sample_id"`
	Status             SpeakerSeparationStatus `json:"status"`
	Speakers           map[string]Speaker      `json:"speakers"`
	SelectedSpeakerIds []string                `json:"selected_speaker_ids"`
}

type Sample struct {
	SampleId                string             `json:"sample_id"`
	FileName                string             `json:"file_name"`
	MimeType                string             `json:"mime_type"`
	SizeBytes               int                `json:"size_bytes"`
	Hash                    string             `json:"hash"`
	DurationSecs            float64            `json:"duration_secs"`
	RemoveBackgroundNoise   bool               `json:"remove_background_noise"`
	HasIsolatedAudio        bool               `json:"has_isolated_audio"`
	HasIsolatedAudioPreview bool               `json:"has_isolated_audio_preview"`
	SpeakerSeparation       *SpeakerSeparation `json:"speaker_separation"`
	TrimStart               int                `json:"trim_start"`
	TrimEnd                 int                `json:"trim_end"`
}

type Recording struct {
	RecordingId    string `json:"recording_id"`
	MimeType       string `json:"mime_type"`
	SizeBytes      int    `json:"size_bytes"`
	UploadDateUnix int    `json:"upload_date_unix"`
	Transcription  string `json:"transcription"`
}

type VerificationAttempt struct {
	Text                string     `json:"text"`
	DateUnix            int        `json:"date_unix"`
	Accepted            bool       `json:"accepted"`
	Similarity          float64    `json:"similarity"`
	LevenshteinDistance float64    `json:"levenshtein_distance"`
	Recording           *Recording `json:"recording"`
}

type ManualVerificationFile struct {
	FileId         string `json:"file_id"`
	FileName       string `json:"file_name"`
	MimeType       string `json:"mime_type"`
	SizeBytes      int    `json:"size_bytes"`
	UploadDateUnix int    `json:"upload_date_unix"`
}

type ManualVerification struct {
	ExtraText       string                   `json:"extra_text"`
	RequestTimeUnix int                      `json:"request_time_unix"`
	Files           []ManualVerificationFile `json:"files"`
}

type FineTuning struct {
	IsAllowedToFineTune                    bool                  `json:"is_allowed_to_fine_tune"`
	State                                  map[string]string     `json:"state"`
	VerificationFailures                   []string              `json:"verification_failures"`
	VerificationAttemptsCount              int                   `json:"verification_attempts_count"`
	ManualVerificationRequested            bool                  `json:"manual_verification_requested"`
	Language                               string                `json:"language"`
	Progress                               map[string]float64    `json:"progress"`
	Message                                map[string]string     `json:"message"`
	DatasetDurationSeconds                 float64               `json:"dataset_duration_seconds"`
	VerificationAttempts                   []VerificationAttempt `json:"verification_attempts"`
	SliceIds                               []string              `json:"slice_ids"`
	ManualVerification                     *ManualVerification   `json:"manual_verification"`
	MaxVerificationAttempts                int                   `json:"max_verification_attempts"`
	NextMaxVerificationAttemptsResetUnixMs int                   `json:"next_max_verification_attempts_reset_unix_ms"`
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	Speed           float64 `json:"speed"`
	CloneStrength   float64 `json:"clone_strength"`
}

type VoiceSharingModerationCheck struct {
	DateCheckedUnix  int       `json:"date_checked_unix"`
	NameValue        string    `json:"name_value"`
	NameCheck        bool      `json:"name_check"`
	DescriptionValue string    `json:"description_value"`
	DescriptionCheck bool      `json:"description_check"`
	SampleIds        []string  `json:"sample_ids"`
	SampleChecks     []float64 `json:"sample_checks"`
	CaptchaIds       []string  `json:"captcha_ids"`
	CaptchaChecks    []float64 `json:"captcha_checks"`
}

type ReaderResource struct {
	ResourceType ReaderResourceType `json:"resource_type"`
	ResourceId   string             `json:"resource_id"`
}

type VoiceSharing struct {
	Status                  VoiceSharingStatus           `json:"status"`
	HistoryItemSampleId     string                       `json:"history_item_sample_id"`
	DateUnix                int                          `json:"date_unix"`
	WhitelistedEmails       []string                     `json:"whitelisted_emails"`
	PublicOwnerId           string                       `json:"public_owner_id"`
	OriginalVoiceId         string                       `json:"original_voice_id"`
	FinancialRewardsEnabled bool                         `json:"financial_rewards_enabled"`
	FreeUsersAllowed        bool                         `json:"free_users_allowed"`
	LiveModerationEnabled   bool                         `json:"live_moderation_enabled"`
	Rate                    float64                      `json:"rate"`
	FiatRate                float64                      `json:"fiat_rate"`
	NoticePeriod            int                          `json:"notice_period"`
	DisableAtUnix           int                          `json:"disable_at_unix"`
	VoiceMixingAllowed      bool                         `json:"voice_mixing_allowed"`
	Featured                bool                         `json:"featured"`
	Category                VoiceCategory                `json:"category"`
	ReaderAppEnabled        bool                         `json:"reader_app_enabled"`
	ImageUrl                string                       `json:"image_url"`
	BanReason               string                       `json:"ban_reason"`
	LikedByCount            int                          `json:"liked_by_count"`
	ClonedByCount           int                          `json:"cloned_by_count"`
	Name                    string                       `json:"name"`
	Description             string                       `json:"description"`
	Labels                  map[string]string            `json:"labels"`
	ReviewStatus            ReviewStatus                 `json:"review_status"`
	ReviewMessage           string                       `json:"review_message"`
	EnabledInLibrary        bool                         `json:"enabled_in_library"`
	InstagramUsername       string                       `json:"instagram_username"`
	TwitterUsername         string                       `json:"twitter_username"`
	YoutubeUsername         string                       `json:"youtube_username"`
	TiktokUsername          string                       `json:"tiktok_username"`
	ModerationCheck         *VoiceSharingModerationCheck `json:"moderation_check"`
	ReaderRestrictedOn      []ReaderResource             `json:"reader_restricted_on"`
}

type VerifiedLanguage struct {
	Language   string `json:"language"`
	ModelId    string `json:"model_id"`
	Accent     string `json:"accent"`
	Locale     string `json:"locale"`
	PreviewUrl string `json:"preview_url"`
}

type VoiceVerification struct {
	RequiresVerification      bool                  `json:"requires_verification"`
	IsVerified                bool                  `json:"is_verified"`
	VerificationFailures      []string              `json:"verification_failures"`
	VerificationAttemptsCount int                   `json:"verification_attempts_count"`
	Language                  string                `json:"language"`
	VerificationAttempts      []VerificationAttempt `json:"verification_attempts"`
}

type Voice struct {
	VoiceId                 string             `json:"voice_id"`
	Name                    string             `json:"name"`
	Samples                 []Sample           `json:"samples"`
	Category                VoiceCategory      `json:"category"`
	FineTuning              *FineTuning        `json:"fine_tuning"`
	Labels                  map[string]string  `json:"labels"`
	Description             string             `json:"description"`
	PreviewUrl              string             `json:"preview_url"`
	AvailableForTiers       []string           `json:"available_for_tiers"`
	Settings                *VoiceSettings     `json:"settings"`
	Sharing                 *VoiceSharing      `json:"sharing"`
	HighQualityBaseModelIds []string           `json:"high_quality_base_model_ids"`
	VerifiedLanguages       []VerifiedLanguage `json:"verified_languages"`
	CollectionIds           []string           `json:"collection_ids"`
	SafetyControl           SafetyControl      `json:"safety_control"`
	VoiceVerification       *VoiceVerification `json:"voice_verification"`
	PermissionOnResource    string             `json:"permission_on_resource"`
	IsOwner                 bool               `json:"is_owner"`
	IsLegacy                bool               `json:"is_legacy"`
	IsMixed                 bool               `json:"is_mixed"`
	FavoritedAtUnix         int                `json:"favorited_at_unix"`
	CreatedAtUnix           int                `json:"created_at_unix"`
	IsBookmarked            bool               `json:"is_bookmarked"`
	RecordingQuality        RecordingQuality   `json:"recording_quality"`
	LabellingStatus         LabellingStatus    `json:"labelling_status"`
	RecordingQualityReason  string             `json:"recording_quality_reason"`
}

type Voices []Voice

type VoiceSearchResult struct {
	Voices        Voices `json:"voices"`
	HasMore       bool   `json:"has_more"`
	TotalCount    int    `json:"total_count"`
	NextPageToken string `json:"next_page_token"`
}

type VoiceFilterFunc func(*Voice) bool

func (v Voices) All() iter.Seq[Voice] {
	return func(yield func(Voice) bool) {
		for _, voice := range v {
			if !yield(voice) {
				return
			}
		}
	}
}

func (v Voices) AllWith(filters ...VoiceFilterFunc) iter.Seq[Voice] {
	return func(yield func(Voice) bool) {
		for _, voice := range v {
			if !matchAllVoice(&voice, filters) {
				continue
			}
			if !yield(voice) {
				return
			}
		}
	}
}

func matchAllVoice(voice *Voice, filters []VoiceFilterFunc) bool {
	for _, filter := range filters {
		if !filter(voice) {
			return false
		}
	}
	return true
}
