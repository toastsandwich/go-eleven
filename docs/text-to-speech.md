# Text to speech

Convert text to a full audio file with `POST /v1/text-to-speech/{voice_id}`.

For chunked audio as it is generated, see [text to speech stream](text-to-speech-stream.md).

[Back to docs](README.md)

## Convert

`TextToSpeech` returns the complete audio as bytes. `text` and `voiceId` are required.

```go
audio, err := client.TextToSpeech(
	"Hello from go eleven.",
	voiceId,
	eleven.WithModelId("eleven_multilingual_v2"),
	eleven.WithOutputFormat(eleven.MP3_44100_128),
	eleven.WithVoiceSettings(eleven.VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Speed:           1,
	}),
)
if err != nil {
	log.Fatal(err)
}
if err := os.WriteFile("speech.mp3", audio, 0644); err != nil {
	log.Fatal(err)
}
```

## Options

Convert and [stream](text-to-speech-stream.md) share these options.

| Option | Where | Notes |
| --- | --- | --- |
| `WithoutEnableLogging` | query | Zero retention. Enterprise only |
| `WithOptimizeStreamingLatency` | query | `StreamingLatencyDefault` … `StreamingLatencyMaxNoNormalizer` |
| `WithOutputFormat` | query | Default `mp3_44100_128` |
| `WithModelId` | body | Default `eleven_multilingual_v2` |
| `WithLanguageCode` | body | ISO 639-1. Not supported on multilingual v2 |
| `WithVoiceSettings` | body | Per-request stability, similarity, speed, style |
| `WithPronunciationDictionaryLocators` | body | Up to 3 locators |
| `WithSeed` | body | Best-effort determinism, `0`–`4294967295` |
| `WithPreviousText` / `WithNextText` | body | Continuity across clips |
| `WithPreviousRequestIds` / `WithNextRequestIds` | body | Up to 3 IDs. Takes priority over previous/next text |
| `WithUsePvcAsIvc` | body | Use IVC of a professional voice |
| `WithApplyTextNormalization` | body | `auto`, `on`, `off` |
| `WithApplyLanguageTextNormalization` | body | Currently Japanese only. Adds latency |
| `WithAudioEffects` | body | Filter, reverb, noise, pan |

Output format examples: `MP3_44100_128`, `MP3_22050_32`, `PCM_16000`, `ULAW8000`, `OPUS_48000_128`. `mp3_44100_192` needs Creator or above. PCM/WAV at 44.1 kHz need Pro or above.

## Context

```go
audio, err := client.TextToSpeechWithContext(ctx, text, voiceId)
```
