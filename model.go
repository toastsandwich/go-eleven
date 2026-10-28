package eleven

import (
	"iter"
)

type Language struct {
	LanguageId string `json:"language_id"`
	Name       string `json:"name"`
}

type ModelRates struct {
	CharacterCostMultiplier float64 `json:"character_cost_multiplier"`
	CostDiscountMultiplier  float64 `json:"cost_discount_multiplier"`
}

type Model struct {
	ModelId                            string     `json:"model_id"`
	Name                               string     `json:"name"`
	CanBeFinetuned                     bool       `json:"can_be_finetuned"`
	CanDoTextToSpeech                  bool       `json:"can_do_text_to_speech"`
	CanDoVoiceConversion               bool       `json:"can_do_voice_conversion"`
	CanUseStyle                        bool       `json:"can_use_style"`
	CanUseSpeakerBoost                 bool       `json:"can_use_speaker_boost"`
	ServesProVoices                    bool       `json:"serves_pro_voices"`
	TokenCostFactor                    float64    `json:"token_cost_factor"`
	Description                        string     `json:"description"`
	RequiresAlphaAccess                bool       `json:"requires_alpha_access"`
	MaximumTextLengthPerRequest        int        `json:"maximum_text_length_per_request"`
	Languages                          []Language `json:"languages"`
	ModelRates                         ModelRates `json:"model_rates"`
	ConcurrencyGroup                   string     `json:"concurrency_group"`
	MaxCharactersRequestFreeUser       int        `json:"max_characters_request_free_user"`
	MaxCharactersRequestSubscribedUser int        `json:"max_characters_request_subscribed_user"`
}

type FilterFunc func(*Model) bool

type Models []Model

func (m Models) All() iter.Seq[Model] {
	return func(yield func(Model) bool) {
		for _, model := range m {
			if !yield(model) {
				return
			}
		}
	}
}

func (m Models) AllWith(filters ...FilterFunc) iter.Seq[Model] {
	return func(yield func(Model) bool) {
		for _, model := range m {
			if !matchAll(&model, filters) {
				continue
			}
			if !yield(model) {
				return
			}
		}
	}
}

func matchAll(model *Model, filters []FilterFunc) bool {
	for _, filter := range filters {
		if !filter(model) {
			return false
		}
	}
	return true
}
