# Text to speech stream

Stream audio as it is generated with `POST /v1/text-to-speech/{voice_id}/stream`.

For a full file in one response, see [text to speech](text-to-speech.md).

[Back to docs](README.md)

## Stream

`TextToSpeechStream` returns an `io.ReadCloser`. Read audio as it arrives, then close the stream. `text` and `voiceId` are required.

```go
stream, err := client.TextToSpeechStream(
	"Hello from go eleven.",
	voiceId,
	eleven.WithModelId("eleven_multilingual_v2"),
	eleven.WithOutputFormat(eleven.MP3_44100_128),
)
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

f, err := os.Create("speech.mp3")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

if _, err := io.Copy(f, stream); err != nil {
	log.Fatal(err)
}
```

The stream client does not use the 10s HTTP timeout so the body can keep flowing. Cancel with context.

Stream does not support WAV output formats.

## Options

Convert and stream share the same options. See [text to speech options](text-to-speech.md#options).

## Context

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

stream, err := client.TextToSpeechStreamWithContext(ctx, text, voiceId)
```
