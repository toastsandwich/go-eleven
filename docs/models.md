# Models

List available ElevenLabs models with `GET /v1/models`.

[Back to docs](README.md)

## List models

```go
models, err := client.ListModels()
if err != nil {
	log.Fatal(err)
}

for model := range models.All() {
	fmt.Println(model.Name, model.ModelId)
}
```

## Filter models

`AllWith` takes one or more filters. A model is yielded only if every filter returns true.

```go
for model := range models.AllWith(func(m *eleven.Model) bool {
	return m.CanDoTextToSpeech
}) {
	fmt.Println(model.Name, model.ModelId)
}
```

## Context

```go
models, err := client.ListModelsWithContext(ctx)
```
