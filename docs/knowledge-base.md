# Knowledge base

Create, list, update, and delete workspace documents with `/v1/convai/knowledge-base`. Attach a document to an agent with `AgentPrompt.KnowledgeBase`. See [agents](agents.md).

[Back to docs](README.md)

## Create document

Create from text, a URL, or an uploaded file (PDF, TXT, DOCX, HTML, EPUB, Markdown, up to 20MB).

```go
doc, err := client.CreateKnowledgeBaseText(
	"Office hours are 9–17 UTC.",
	eleven.WithKnowledgeName("hours"),
)
if err != nil {
	log.Fatal(err)
}

fromURL, err := client.CreateKnowledgeBaseURL(
	"https://example.com/docs",
	eleven.WithKnowledgeURLName("product-docs"),
	eleven.WithKnowledgeAutoSync(),
)

file, err := os.Open("faq.txt")
if err != nil {
	log.Fatal(err)
}
defer file.Close()

fromFile, err := client.CreateKnowledgeBaseFile(
	file,
	eleven.WithKnowledgeFileName("faq"),
	eleven.WithKnowledgeFilename("faq.txt"),
)
```

Create always inserts a new document. Persist `id` if you need the same document after a restart.

## List documents

```go
result, err := client.ListKnowledgeBase(
	eleven.WithKnowledgeSearch("hours"),
	eleven.WithKnowledgePageSize(30),
	eleven.WithKnowledgeTypes(eleven.KnowledgeTypeText, eleven.KnowledgeTypeURL),
	eleven.WithKnowledgeSort(eleven.KnowledgeSortName),
	eleven.WithKnowledgeSortDirection(eleven.SortAsc),
)
if err != nil {
	log.Fatal(err)
}

for document := range result.Documents.All() {
	fmt.Println(document.ID, document.Name, document.Type)
}
```

`Documents` also has `AllWith` for extra local filters after the API response.

Paginate with `has_more` and `next_cursor`.

```go
if result.HasMore {
	result, err = client.ListKnowledgeBase(
		eleven.WithKnowledgePageSize(30),
		eleven.WithKnowledgeCursor(result.NextCursor),
	)
}
```

## Get, content, update, RAG, delete

```go
document, err := client.GetKnowledgeBase(documentationID)
content, err := client.GetKnowledgeBaseContent(documentationID)

document, err = client.UpdateKnowledgeBase(
	documentationID,
	eleven.WithUpdateKnowledgeName("hours v2"),
	eleven.WithUpdateKnowledgeContent("Office hours are 10–18 UTC."),
)

index, err := client.ComputeRAGIndex(documentationID, eleven.EmbeddingModelE5Mistral7BInstruct)
chunk, err := client.GetKnowledgeBaseChunk(documentationID, chunkID)

dependents, err := client.GetKnowledgeBaseDependentAgents(documentationID)
err = client.DeleteKnowledgeBase(documentationID, eleven.WithForceDeleteKnowledge())
err = client.DeleteKnowledgeBase(documentationID)
```

`GetKnowledgeBaseContent` returns the extracted document body. `ComputeRAGIndex` starts indexing when the document is not indexed yet, and returns the current status otherwise. Delete fails while agents still reference the document unless you pass `WithForceDeleteKnowledge`.

## Attach to an agent

```go
_, err = client.UpdateAgent(agentID, eleven.WithUpdateConversationConfig(eleven.ConversationalConfig{
	Agent: &eleven.AgentConfig{
		Prompt: &eleven.AgentPrompt{
			KnowledgeBase: []eleven.AgentKnowledgeBaseLocator{{
				Type: eleven.KnowledgeTypeText,
				Name: doc.Name,
				ID:   doc.ID,
			}},
			RAG: json.RawMessage(`{"enabled":true,"embedding_model":"e5_mistral_7b_instruct"}`),
		},
	},
}))
```

## Options

### List

| Option | Query | Notes |
| --- | --- | --- |
| `WithKnowledgePageSize` | `page_size` | Max 100, default 30 |
| `WithKnowledgeSearch` | `search` | Name prefix |
| `WithKnowledgeCreatedByUserID` | `created_by_user_id` | `@me` is the authenticated user |
| `WithKnowledgeTypes` | `types` | Repeated; `file`, `url`, `text`, `folder` |
| `WithKnowledgeParentFolderID` | `parent_folder_id` | Direct children |
| `WithKnowledgeAncestorFolderID` | `ancestor_folder_id` | Descendants |
| `WithKnowledgeFoldersFirst` | `folders_first=true` | |
| `WithKnowledgeSort` | `sort_by` | `name`, `created_at`, `updated_at`, `size` |
| `WithKnowledgeSortDirection` | `sort_direction` | `SortAsc` / `SortDesc` |
| `WithKnowledgeCursor` | `cursor` | Next page from `NextCursor` |

### Create text

| Option | Body | Notes |
| --- | --- | --- |
| `WithKnowledgeName` | `name` | |
| `WithKnowledgeTextParentFolderID` | `parent_folder_id` | |

### Create URL

| Option | Body | Notes |
| --- | --- | --- |
| `WithKnowledgeURLName` | `name` | |
| `WithKnowledgeURLParentFolderID` | `parent_folder_id` | |
| `WithKnowledgeAutoSync` | `enable_auto_sync=true` | |
| `WithKnowledgeAutoRemove` | `auto_remove=true` | Used with auto-sync |
| `WithKnowledgeMinFrequencyDays` | `minimum_frequency_days` | Default 7 |

### Create file

| Option | Form | Notes |
| --- | --- | --- |
| `WithKnowledgeFileName` | `name` | |
| `WithKnowledgeFileParentFolderID` | `parent_folder_id` | |
| `WithKnowledgeFilename` | file part name | Defaults to the `os.File` name, otherwise `document` |

### Get / update / delete / RAG

| Option | Where | Notes |
| --- | --- | --- |
| `WithKnowledgeAgentID` | query | Scope get to an agent |
| `WithUpdateKnowledgeName` | body | |
| `WithUpdateKnowledgeContent` | body | Text docs, URL without auto-sync, unsynced files |
| `WithForceDeleteKnowledge` | query | Delete even if agents use the document |
| `WithKnowledgeDependentCursor` | query | Next page from `NextCursor` |
| `WithKnowledgeDependentPageSize` | query | Max 100, default 30 |
| `WithKnowledgeEmbeddingModel` | query | Chunk lookup |

## Context

```go
list, err := client.ListKnowledgeBaseWithContext(ctx, eleven.WithKnowledgeSearch("hours"))
doc, err := client.CreateKnowledgeBaseTextWithContext(ctx, "Office hours are 9–17 UTC.", eleven.WithKnowledgeName("hours"))
document, err := client.GetKnowledgeBaseWithContext(ctx, documentationID)
err = client.DeleteKnowledgeBaseWithContext(ctx, documentationID)
```
