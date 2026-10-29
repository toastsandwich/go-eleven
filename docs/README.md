# go-eleven

Go SDK for the [ElevenLabs](https://elevenlabs.io) API.

## Install

```bash
go get github.com/toastsandwich/go-eleven
```

## Setup

Create an API key in the ElevenLabs dashboard. Keys start with `sk_`.

```go
package main

import (
	"os"

	"github.com/toastsandwich/go-eleven"
)

func main() {
	client := eleven.NewClient(os.Getenv("ELEVEN_API_KEY"))
}
```

`NewClient` panics if the key is empty. Every method also has a `WithContext` variant.

## APIs

- [Models](models.md) — list models and filter with `All` / `AllWith`
- [Voices](voices.md) — search, filter, and paginate voices
- [Text to speech](text-to-speech.md) — convert text to a full audio file
- [Text to speech stream](text-to-speech-stream.md) — stream audio as it is generated
- [Webhooks](webhooks.md) — manage workspace webhooks and verify deliveries
