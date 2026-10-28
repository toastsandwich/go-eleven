# go-eleven

Go SDK for the [ElevenLabs](https://elevenlabs.io) API.

```bash
go get github.com/toastsandwich/go-eleven
```

```go
client := eleven.NewClient(os.Getenv("ELEVEN_API_KEY"))
```

Supported APIs:

- [Models](docs/models.md)
- [Voices](docs/voices.md)
- [Text to speech](docs/text-to-speech.md)
- [Text to speech stream](docs/text-to-speech-stream.md)

See [docs/README.md](docs/README.md) for setup and usage.
