# Speech to speech

Transform audio from one voice to another with `POST /v1/voice-changer/{voice_id}`. Emotion, timing, and delivery of the source recording are preserved.

For chunked audio as it is generated, see [speech to speech stream](speech-to-speech-stream.md).

[Back to docs](README.md)

## Convert

`SpeechToSpeech` uploads source audio as multipart form data and returns the complete converted audio as bytes. `audio` and `voiceId` are required.

```go
f, err := os.Open("source.mp3")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

audio, err := client.SpeechToSpeech(
	f,
	voiceId,
	eleven.WithSpeechToSpeechModelId("eleven_multilingual_sts_v2"),
	eleven.WithSpeechToSpeechOutputFormat(eleven.MP3_44100_128),
	eleven.WithSpeechToSpeechVoiceSettings(eleven.VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Speed:           1,
	}),
)
if err != nil {
	log.Fatal(err)
}
if err := os.WriteFile("converted.mp3", audio, 0644); err != nil {
	log.Fatal(err)
}
```

The HTTP timeout is disabled so conversion of longer clips can finish. Cancel with context.

Source audio is at most 5 minutes and 50 MB per request.

## Options

Convert and [stream](speech-to-speech-stream.md) share these options.

| Option | Where | Notes |
| --- | --- | --- |
| `WithoutSpeechToSpeechEnableLogging` | query | Zero retention. Enterprise only |
| `WithSpeechToSpeechOptimizeStreamingLatency` | query | `StreamingLatencyDefault` … `StreamingLatencyMaxNoNormalizer` |
| `WithSpeechToSpeechOutputFormat` | query | Default `mp3_44100_128` |
| `WithSpeechToSpeechModelId` | form | Default `eleven_english_sts_v2`. Prefer `eleven_multilingual_sts_v2` |
| `WithSpeechToSpeechVoiceSettings` | form | JSON-encoded per-request stability, similarity, speed, style |
| `WithSpeechToSpeechSeed` | form | Best-effort determinism, `0`–`4294967295` |
| `WithRemoveBackgroundNoise` | form | Isolate the source voice before conversion |
| `WithFileFormat` | form | `FileFormatOther` or `FileFormatPCMS16LE16` |
| `WithAudioFilename` | form | Multipart filename. Defaults to an `*os.File` base name, or `audio` |

The model must support speech to speech (`can_do_voice_conversion`). Output format examples: `MP3_44100_128`, `PCM_16000`, `ULAW8000`, `OPUS_48000_64`. `mp3_44100_192` needs Creator or above. PCM/WAV at 44.1 kHz need Pro or above.

`FileFormatPCMS16LE16` is 16-bit PCM, 16 kHz, mono, little-endian. That skips decoding and lowers latency.

## Context

```go
audio, err := client.SpeechToSpeechWithContext(ctx, source, voiceId)
```
