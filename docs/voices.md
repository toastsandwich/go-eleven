# Voices

Search, filter, and paginate voices with `GET /v2/voices`.

[Back to docs](README.md)

## Search voices

```go
result, err := client.SearchVoices(
	eleven.WithPageSize(20),
	eleven.WithSearch("narrator"),
	eleven.WithVoiceType(eleven.VoiceTypeDefault),
	eleven.WithGender("male"),
	eleven.WithSort(eleven.VoiceSortName),
	eleven.WithSortDirection(eleven.SortAsc),
)
if err != nil {
	log.Fatal(err)
}

for voice := range result.Voices.All() {
	fmt.Println(voice.Name, voice.VoiceId)
}
```

`Voices` also has `AllWith` for extra local filters after the API response.

## Pagination

Paginate with `has_more` and `next_page_token`. `total_count` is a live snapshot and can change between pages.

```go
if result.HasMore {
	result, err = client.SearchVoices(
		eleven.WithPageSize(20),
		eleven.WithNextPageToken(result.NextPageToken),
	)
}
```

## Options

| Option | Purpose |
| --- | --- |
| `WithSearch` | Match name, description, labels, category |
| `WithPageSize` | Max voices per page (up to 100, default 10) |
| `WithNextPageToken` | Continue from the previous page |
| `WithVoiceType` | `personal`, `community`, `default`, `workspace`, `saved`, … |
| `WithCategory` | `premade`, `cloned`, `generated`, `professional` |
| `WithFineTuningState` | Professional voice clone state |
| `WithCollectionId` | Voices in a collection |
| `WithGender`, `WithAge`, `WithAccent` | Label filters |
| `WithLanguages`, `WithUseCases` | Match any of the given values |
| `WithVoiceIds` | Look up specific voice IDs |
| `WithHighQuality` | Studio-quality voices only |
| `WithIncludeCustomRates` | Include custom sharing rates |
| `WithIncludeLiveModerated` | Include live-moderated voices |
| `WithIncludeTotalCount` | Include `total_count` (extra cost) |
| `WithMinNoticePeriodDays` | Minimum sharing notice period |
| `WithSort` | `created_at_unix` or `name` |
| `WithSortDirection` | `asc` or `desc` |

## Context

```go
result, err := client.SearchVoicesWithContext(ctx, eleven.WithSearch("news"))
```
