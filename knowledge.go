package eleven

import (
	"encoding/json"
	"iter"
)

type KnowledgeType string

const (
	KnowledgeTypeFile   KnowledgeType = "file"
	KnowledgeTypeURL    KnowledgeType = "url"
	KnowledgeTypeText   KnowledgeType = "text"
	KnowledgeTypeFolder KnowledgeType = "folder"
)

func (t KnowledgeType) String() string { return string(t) }

type KnowledgeSort string

const (
	KnowledgeSortName      KnowledgeSort = "name"
	KnowledgeSortCreatedAt KnowledgeSort = "created_at"
	KnowledgeSortUpdatedAt KnowledgeSort = "updated_at"
	KnowledgeSortSize      KnowledgeSort = "size"
)

func (s KnowledgeSort) String() string { return string(s) }

type KnowledgeUsage string

const (
	KnowledgeUsagePrompt KnowledgeUsage = "prompt"
	KnowledgeUsageAuto   KnowledgeUsage = "auto"
)

func (u KnowledgeUsage) String() string { return string(u) }

type EmbeddingModel string

const (
	EmbeddingModelE5Mistral7BInstruct         EmbeddingModel = "e5_mistral_7b_instruct"
	EmbeddingModelMultilingualE5LargeInstruct EmbeddingModel = "multilingual_e5_large_instruct"
)

func (m EmbeddingModel) String() string { return string(m) }

type RAGIndexStatus string

const (
	RAGIndexNew               RAGIndexStatus = "new"
	RAGIndexCreated           RAGIndexStatus = "created"
	RAGIndexProcessing        RAGIndexStatus = "processing"
	RAGIndexFailed            RAGIndexStatus = "failed"
	RAGIndexSucceeded         RAGIndexStatus = "succeeded"
	RAGIndexLimitExceeded     RAGIndexStatus = "rag_limit_exceeded"
	RAGIndexDocumentTooSmall  RAGIndexStatus = "document_too_small"
	RAGIndexCannotIndexFolder RAGIndexStatus = "cannot_index_folder"
)

func (s RAGIndexStatus) String() string { return string(s) }

type KnowledgeMetadata struct {
	CreatedAtUnixSecs     int64 `json:"created_at_unix_secs"`
	LastUpdatedAtUnixSecs int64 `json:"last_updated_at_unix_secs"`
	SizeBytes             int64 `json:"size_bytes"`
}

type KnowledgeFolderPathSegment struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type KnowledgeDocument struct {
	Type             KnowledgeType                `json:"type"`
	ID               string                       `json:"id"`
	Name             string                       `json:"name"`
	AccessInfo       ResourceAccessInfo           `json:"access_info"`
	Metadata         KnowledgeMetadata            `json:"metadata"`
	SupportedUsages  []KnowledgeUsage             `json:"supported_usages"`
	URL              string                       `json:"url,omitempty"`
	Filename         string                       `json:"filename,omitempty"`
	ExtractedHTML    string                       `json:"extracted_inner_html,omitempty"`
	ContentFormat    string                       `json:"content_format,omitempty"`
	FolderParentID   string                       `json:"folder_parent_id,omitempty"`
	FolderPath       []KnowledgeFolderPathSegment `json:"folder_path,omitempty"`
	ChildrenCount    int                          `json:"children_count,omitempty"`
	DocumentCount    int                          `json:"document_count,omitempty"`
	AutoSyncInfo     json.RawMessage              `json:"auto_sync_info,omitempty"`
	ExternalSyncInfo json.RawMessage              `json:"external_sync_info,omitempty"`
}

type KnowledgeDocuments []KnowledgeDocument

type KnowledgeFilterFunc func(*KnowledgeDocument) bool

func (d KnowledgeDocuments) All() iter.Seq[KnowledgeDocument] {
	return func(yield func(KnowledgeDocument) bool) {
		for _, document := range d {
			if !yield(document) {
				return
			}
		}
	}
}

func (d KnowledgeDocuments) AllWith(filters ...KnowledgeFilterFunc) iter.Seq[KnowledgeDocument] {
	return func(yield func(KnowledgeDocument) bool) {
		for _, document := range d {
			if !matchAllKnowledge(&document, filters) {
				continue
			}
			if !yield(document) {
				return
			}
		}
	}
}

func matchAllKnowledge(document *KnowledgeDocument, filters []KnowledgeFilterFunc) bool {
	for _, filter := range filters {
		if !filter(document) {
			return false
		}
	}
	return true
}

type KnowledgeList struct {
	Documents  KnowledgeDocuments `json:"documents"`
	HasMore    bool               `json:"has_more"`
	NextCursor string             `json:"next_cursor"`
}

type CreateKnowledgeResult struct {
	ID         string                       `json:"id"`
	Name       string                       `json:"name"`
	FolderPath []KnowledgeFolderPathSegment `json:"folder_path,omitempty"`
}

type KnowledgeContent struct {
	Content string `json:"content,omitempty"`
}

func decodeKnowledgeContent(body []byte) KnowledgeContent {
	var asString string
	if err := json.Unmarshal(body, &asString); err == nil {
		return KnowledgeContent{Content: asString}
	}
	var asObject KnowledgeContent
	if err := json.Unmarshal(body, &asObject); err == nil && asObject.Content != "" {
		return asObject
	}
	return KnowledgeContent{Content: string(body)}
}

type KnowledgeChunk struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type RAGDocumentIndexUsage struct {
	UsedBytes int64 `json:"used_bytes"`
}

type RAGIndex struct {
	ID                      string                `json:"id"`
	Model                   EmbeddingModel        `json:"model"`
	Status                  RAGIndexStatus        `json:"status"`
	ProgressPercentage      float64               `json:"progress_percentage"`
	DocumentModelIndexUsage RAGDocumentIndexUsage `json:"document_model_index_usage"`
}
