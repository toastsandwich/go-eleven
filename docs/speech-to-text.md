# Speech to text

Transcribe audio or video with `POST /v1/speech-to-text`. Retrieve a stored transcript with `GET /v1/speech-to-text/transcripts/{transcription_id}`. Delete one with `DELETE /v1/speech-to-text/transcripts/{transcription_id}`.

This is the batch Scribe API. Provide a file or a `source_url`, not both.

[Back to docs](README.md)

## Convert

`SpeechToText` uploads source audio as multipart form data and returns a `SpeechToTextResult`. `model_id` defaults to `scribe_v2`. Use `scribe_v2_medical` for medical and clinical audio.

```go
f, err := os.Open("source.mp3")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

result, err := client.SpeechToText(
	f,
	eleven.WithSpeechToTextModelId("scribe_v2"),
	eleven.WithSpeechToTextLanguageCode("eng"),
	eleven.WithDiarize(),
	eleven.WithTagAudioEvents(),
	eleven.WithKeyterms("ElevenLabs", "Scribe"),
)
if err != nil {
	log.Fatal(err)
}
fmt.Println(result.Text, result.LanguageCode, result.TranscriptionID)
```

The HTTP timeout is disabled so transcription of longer clips can finish. Cancel with context.

Audio is at least 100ms. Files are at most 5GB. All major audio and video formats are supported.

## Source URL

Skip the file upload and transcribe a hosted file, YouTube URL, TikTok URL, or other video host.

```go
result, err := client.SpeechToText(
	nil,
	eleven.WithSourceURL("https://example.com/meeting.mp3"),
	eleven.WithSpeechToTextModelId("scribe_v2"),
)
```

## Result shapes

A synchronous request returns a single transcript on `SpeechToTextResult` (`Text`, `Words`, `LanguageCode`, and so on).

`WithUseMultiChannel` with style `separate` fills `Transcripts`, one per channel. `All` / `AllWith` iterate those transcripts, or the single transcript when there is only one.

`WithSpeechToTextWebhook` returns early. `IsWebhook` is true, and `RequestID` / `TranscriptionID` identify the job. The finished transcript is delivered to configured `speech_to_text` webhooks. Decode it with `SpeechToTextTranscriptionData.Result`.

## Get and delete

```go
result, err := client.GetTranscript(transcriptionID)
if err != nil {
	log.Fatal(err)
}

if err := client.DeleteTranscript(transcriptionID); err != nil {
	log.Fatal(err)
}
```

## Options

| Option | Where | Notes |
| --- | --- | --- |
| `WithSpeechToTextToken` | query | Single-use `batch_scribe` token |
| `WithoutSpeechToTextEnableLogging` | query | Zero retention. Enterprise only |
| `WithSpeechToTextModelId` | form | Default `scribe_v2` |
| `WithSpeechToTextLanguageCode` | form | ISO 639-1 or 639-3. Detected when omitted |
| `WithTranscriptEdit` | form | Natural-language edit. Extra 30% surcharge |
| `WithTagAudioEvents` / `WithoutTagAudioEvents` | form | Tag laughter, footsteps, and similar events |
| `WithNumSpeakers` | form | Max 32 |
| `WithTimestampsGranularity` | form | `none`, `word`, `character` |
| `WithDiarize` | form | Annotate speakers |
| `WithDiarizationThreshold` | form | Only with diarize and without `num_speakers` |
| `WithAdditionalFormats` | form | `docx`, `html`, `pdf`, `segmented_json`, `srt`, `txt` |
| `WithSpeechToTextFileFormat` | form | `FileFormatOther` or `FileFormatPCMS16LE16` |
| `WithSourceURL` | form | Instead of uploading a file |
| `WithSpeechToTextWebhook` | form | Async; result sent to webhooks |
| `WithSpeechToTextWebhookID` | form | Target one webhook |
| `WithSpeechToTextWebhookMetadata` | form | JSON object, max 16KB, depth 2 |
| `WithSpeechToTextTemperature` | form | `0.0`–`2.0` |
| `WithSpeechToTextSeed` | form | Best-effort determinism, `0`–`2147483647` |
| `WithUseMultiChannel` | form | Up to 5 channels, billed per channel |
| `WithMultichannelOutputStyle` | form | `separate` or `combined` |
| `WithEntityDetection` | form | `all`, a category, or entity types |
| `WithEntityRedaction` | form | Subset of detection |
| `WithEntityRedactionMode` | form | `redacted`, `entity_type`, `enumerated_entity_type` |
| `WithNoVerbatim` | form | Drop fillers. `scribe_v2` only |
| `WithUseSpeakerLibrary` | form | Match known speakers when diarizing |
| `WithDetectSpeakerRoles` | form | `agent` / `customer`. Extra 10% surcharge |
| `WithKeyterms` | form | Up to 1000 terms. Extra 20% surcharge |
| `WithSpeechToTextAudioFilename` | form | Multipart filename. Defaults to an `*os.File` base name, or `audio` |

`FileFormatPCMS16LE16` is 16-bit PCM, 16 kHz, mono, little-endian. That skips decoding and lowers latency.

## Context

```go
result, err := client.SpeechToTextWithContext(ctx, source)
result, err = client.GetTranscriptWithContext(ctx, transcriptionID)
err = client.DeleteTranscriptWithContext(ctx, transcriptionID)
```
