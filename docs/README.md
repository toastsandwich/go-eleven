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
- [Speech to speech](speech-to-speech.md) — convert a recording into another voice
- [Speech to speech stream](speech-to-speech-stream.md) — stream converted audio as it is generated
- [Speech to text](speech-to-text.md) — transcribe audio or video with Scribe
- [Agents](agents.md) — create, list, update, and delete Eleven Agents
- [Conversations](conversations.md) — list calls, fetch transcripts, and mint session tokens
- [Live conversation](live-conversation.md) — talk to an agent over WebSocket
- [Tools](tools.md) — create, list, update, and delete workspace tools
- [Knowledge base](knowledge-base.md) — add documents and compute a RAG index
- [Webhooks](webhooks.md) — manage workspace webhooks and verify deliveries
