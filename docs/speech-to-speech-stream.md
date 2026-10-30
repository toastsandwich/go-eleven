# Speech to speech stream

Stream converted audio as it is generated with `POST /v1/voice-changer/{voice_id}/stream`.

For a full file in one response, see [speech to speech](speech-to-speech.md).

[Back to docs](README.md)

## Stream

`SpeechToSpeechStream` returns an `io.ReadCloser`. Read audio as it arrives, then close the stream. `audio` and `voiceId` are required.

```go
f, err := os.Open("source.mp3")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

stream, err := client.SpeechToSpeechStream(
	f,
	voiceId,
	eleven.WithSpeechToSpeechModelId("eleven_multilingual_sts_v2"),
	eleven.WithSpeechToSpeechOutputFormat(eleven.MP3_44100_128),
)
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

out, err := os.Create("converted.mp3")
if err != nil {
	log.Fatal(err)
}
defer out.Close()

if _, err := io.Copy(out, stream); err != nil {
	log.Fatal(err)
}
```

The stream client does not use the 10s HTTP timeout so the body can keep flowing. Cancel with context.

Stream does not support WAV output formats.

## Options

Convert and stream share the same options. See [speech to speech options](speech-to-speech.md#options).

## Context

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

stream, err := client.SpeechToSpeechStreamWithContext(ctx, source, voiceId)
```
