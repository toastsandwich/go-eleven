package eleven

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestVoicesAllWith(t *testing.T) {
	voices := Voices{
		{Name: "Roger", VoiceId: "1", Category: VoiceCategoryPremade},
		{Name: "Clone", VoiceId: "2", Category: VoiceCategoryCloned},
		{Name: "Studio", VoiceId: "3", Category: VoiceCategoryPremade},
	}

	var names []string
	for voice := range voices.AllWith(func(v *Voice) bool {
		return v.Category == VoiceCategoryPremade
	}) {
		names = append(names, voice.Name)
	}
	if len(names) != 2 || names[0] != "Roger" || names[1] != "Studio" {
		t.Fatalf("names = %v", names)
	}
}

func TestSearchVoicesQuery(t *testing.T) {
	var opt searchVoicesOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithNextPageToken("tok")(&opt)
	WithPageSize(20)(&opt)
	WithSearch("narrator")(&opt)
	WithSort(VoiceSortName)(&opt)
	WithSortDirection(SortAsc)(&opt)
	WithVoiceType(VoiceTypeDefault)(&opt)
	WithCategory(VoiceCategoryPremade)(&opt)
	WithFineTuningState(FineTuningStateFineTuned)(&opt)
	WithCollectionId("col")(&opt)
	WithGender("male")(&opt)
	WithAge("middle_aged")(&opt)
	WithLanguages("en", "es")(&opt)
	WithAccent("american")(&opt)
	WithUseCases("narration")(&opt)
	WithMinNoticePeriodDays(7)(&opt)
	WithIncludeCustomRates(false)(&opt)
	WithIncludeLiveModerated(true)(&opt)
	WithHighQuality(true)(&opt)
	WithIncludeTotalCount(true)(&opt)
	WithVoiceIds("v1", "v2")(&opt)

	q := opt.query()
	assertQuery(t, q, "next_page_token", "tok")
	assertQuery(t, q, "page_size", "20")
	assertQuery(t, q, "search", "narrator")
	assertQuery(t, q, "sort", "name")
	assertQuery(t, q, "sort_direction", "asc")
	assertQuery(t, q, "voice_type", "default")
	assertQuery(t, q, "category", "premade")
	assertQuery(t, q, "fine_tuning_state", "fine_tuned")
	assertQuery(t, q, "collection_id", "col")
	assertQuery(t, q, "gender", "male")
	assertQuery(t, q, "age", "middle_aged")
	assertQuery(t, q, "language", `["en","es"]`)
	assertQuery(t, q, "accent", "american")
	assertQuery(t, q, "use_cases", `["narration"]`)
	assertQuery(t, q, "min_notice_period_days", "7")
	assertQuery(t, q, "include_custom_rates", "false")
	assertQuery(t, q, "include_live_moderated", "true")
	assertQuery(t, q, "high_quality", "true")
	assertQuery(t, q, "include_total_count", "true")
	assertQuery(t, q, "voice_ids", `["v1","v2"]`)
}

func TestSearchVoicesQueryOmitsUnsetBools(t *testing.T) {
	var opt searchVoicesOption
	q := opt.query()
	for _, key := range []string{"include_custom_rates", "include_live_moderated", "high_quality", "include_total_count"} {
		if q.Has(key) {
			t.Fatalf("%s should be omitted when unset", key)
		}
	}
}

func TestVoiceEnumStrings(t *testing.T) {
	cases := []struct {
		got, want string
	}{
		{VoiceCategoryPremade.String(), "premade"},
		{VoiceTypeNonDefault.String(), "non-default"},
		{FineTuningStateNotVerified.String(), "not_verified"},
		{VoiceSortCreatedAtUnix.String(), "created_at_unix"},
		{SortDesc.String(), "desc"},
		{SafetyControlEnterpriseBan.String(), "ENTERPRISE_BAN"},
		{RecordingQualityStudio.String(), "studio"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("got %q want %q", tc.got, tc.want)
		}
	}
}

func TestSearchVoices(t *testing.T) {
	const body = `{"voices":[{"voice_id":"CwhRBWXzGAHq8TQ4Fs17","name":"Roger","category":"premade"}],"has_more":true,"total_count":12,"next_page_token":"next"}`
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v2/voices" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("search") != "roger" || r.URL.Query().Get("page_size") != "5" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})

	result, err := client.SearchVoices(WithSearch("roger"), WithPageSize(5))
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasMore || result.TotalCount != 12 || result.NextPageToken != "next" {
		t.Fatalf("result = %+v", result)
	}
	if len(result.Voices) != 1 || result.Voices[0].Name != "Roger" || result.Voices[0].VoiceId != "CwhRBWXzGAHq8TQ4Fs17" {
		t.Fatalf("voices = %+v", result.Voices)
	}
	if result.Voices[0].Category != VoiceCategoryPremade {
		t.Fatalf("category = %s", result.Voices[0].Category)
	}
}

func TestSearchVoicesError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	})
	_, err := client.SearchVoices()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVoiceSearchResultJSON(t *testing.T) {
	raw := []byte(`{"voices":[{"voice_id":"abc","name":"Aria","labels":{"gender":"female"}}],"has_more":false,"total_count":1}`)
	var result VoiceSearchResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatal(err)
	}
	if result.HasMore || result.Voices[0].Labels["gender"] != "female" {
		t.Fatalf("result = %+v", result)
	}
}

func TestLiveSearchVoices(t *testing.T) {
	result, err := liveClient(t).SearchVoices(
		WithPageSize(5),
		WithVoiceType(VoiceTypeDefault),
		WithIncludeTotalCount(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Voices) == 0 {
		t.Fatal("no voices")
	}
	for voice := range result.Voices.All() {
		if voice.VoiceId == "" || voice.Name == "" {
			t.Fatalf("incomplete voice: %+v", voice)
		}
	}
	t.Logf("voices=%d has_more=%v total_count=%d", len(result.Voices), result.HasMore, result.TotalCount)
}

func assertQuery(t *testing.T, q url.Values, key, want string) {
	t.Helper()
	if got := q.Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
