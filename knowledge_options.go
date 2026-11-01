package eleven

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
)

type listKnowledgeOption struct {
	PageSize         int
	Search           string
	CreatedByUserID  string
	Types            []string
	ParentFolderID   string
	AncestorFolderID string
	FoldersFirst     *bool
	Sort             string
	SortDirection    string
	Cursor           string
}

type ListKnowledgeOption func(*listKnowledgeOption)

// How many documents to return at maximum. Can not exceed 100, defaults to 30.
func WithKnowledgePageSize(size int) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.PageSize = size
	}
}

// If specified, the endpoint returns only documents whose names start with this string.
func WithKnowledgeSearch(search string) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.Search = search
	}
}

// Filter documents by creator user ID. When set, only documents created by this user are returned. Use "@me" to refer to the authenticated user.
func WithKnowledgeCreatedByUserID(id string) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.CreatedByUserID = id
	}
}

// If present, the endpoint returns only documents of the given types. file, url, text, or folder.
func WithKnowledgeTypes(types ...KnowledgeType) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.Types = make([]string, len(types))
		for i, t := range types {
			o.Types[i] = t.String()
		}
	}
}

// If set, the endpoint returns only documents that are direct children of the given folder.
func WithKnowledgeParentFolderID(id string) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.ParentFolderID = id
	}
}

// If set, the endpoint returns only documents that are descendants of the given folder.
func WithKnowledgeAncestorFolderID(id string) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.AncestorFolderID = id
	}
}

// Whether folders should be returned first in the list of documents. Sends folders_first=true.
func WithKnowledgeFoldersFirst() ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		foldersFirst := true
		o.FoldersFirst = &foldersFirst
	}
}

// The field to sort the results by. name, created_at, updated_at, or size.
func WithKnowledgeSort(sort KnowledgeSort) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.Sort = sort.String()
	}
}

// The direction to sort the results. asc or desc.
func WithKnowledgeSortDirection(direction SortDirection) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.SortDirection = direction.String()
	}
}

// Used for fetching the next page. Cursor is returned in the list response as NextCursor.
func WithKnowledgeCursor(cursor string) ListKnowledgeOption {
	return func(o *listKnowledgeOption) {
		o.Cursor = cursor
	}
}

func (o *listKnowledgeOption) query() url.Values {
	q := url.Values{}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.Search != "" {
		q.Set("search", o.Search)
	}
	if o.CreatedByUserID != "" {
		q.Set("created_by_user_id", o.CreatedByUserID)
	}
	for _, t := range o.Types {
		q.Add("types", t)
	}
	if o.ParentFolderID != "" {
		q.Set("parent_folder_id", o.ParentFolderID)
	}
	if o.AncestorFolderID != "" {
		q.Set("ancestor_folder_id", o.AncestorFolderID)
	}
	if o.FoldersFirst != nil {
		q.Set("folders_first", strconv.FormatBool(*o.FoldersFirst))
	}
	if o.Sort != "" {
		q.Set("sort_by", o.Sort)
	}
	if o.SortDirection != "" {
		q.Set("sort_direction", o.SortDirection)
	}
	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	return q
}

type knowledgeNameFolderOption struct {
	Name           string
	ParentFolderID string
}

// A custom, human-readable name for the document.
func withKnowledgeName(name string, o *knowledgeNameFolderOption) {
	o.Name = name
}

// If set, the created document is placed inside the given folder.
func withKnowledgeParentFolderID(id string, o *knowledgeNameFolderOption) {
	o.ParentFolderID = id
}

type createKnowledgeTextOption struct {
	knowledgeNameFolderOption
}

type CreateKnowledgeTextOption func(*createKnowledgeTextOption)

// A custom, human-readable name for the document.
func WithKnowledgeName(name string) CreateKnowledgeTextOption {
	return func(o *createKnowledgeTextOption) {
		withKnowledgeName(name, &o.knowledgeNameFolderOption)
	}
}

// If set, the created document is placed inside the given folder.
func WithKnowledgeTextParentFolderID(id string) CreateKnowledgeTextOption {
	return func(o *createKnowledgeTextOption) {
		withKnowledgeParentFolderID(id, &o.knowledgeNameFolderOption)
	}
}

type createKnowledgeTextRequest struct {
	Text           string `json:"text"`
	Name           string `json:"name,omitempty"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
}

func (o *createKnowledgeTextOption) body(text string) createKnowledgeTextRequest {
	return createKnowledgeTextRequest{
		Text:           text,
		Name:           o.Name,
		ParentFolderID: o.ParentFolderID,
	}
}

type createKnowledgeURLOption struct {
	knowledgeNameFolderOption
	EnableAutoSync       *bool
	AutoRemove           *bool
	MinimumFrequencyDays *int
}

type CreateKnowledgeURLOption func(*createKnowledgeURLOption)

// A custom, human-readable name for the document.
func WithKnowledgeURLName(name string) CreateKnowledgeURLOption {
	return func(o *createKnowledgeURLOption) {
		withKnowledgeName(name, &o.knowledgeNameFolderOption)
	}
}

// If set, the created document is placed inside the given folder.
func WithKnowledgeURLParentFolderID(id string) CreateKnowledgeURLOption {
	return func(o *createKnowledgeURLOption) {
		withKnowledgeParentFolderID(id, &o.knowledgeNameFolderOption)
	}
}

// Whether to enable auto-sync for this URL document. Sends enable_auto_sync=true.
func WithKnowledgeAutoSync() CreateKnowledgeURLOption {
	return func(o *createKnowledgeURLOption) {
		enabled := true
		o.EnableAutoSync = &enabled
	}
}

// Whether to automatically remove the document if the URL becomes unavailable. Only used when auto-sync is enabled.
func WithKnowledgeAutoRemove() CreateKnowledgeURLOption {
	return func(o *createKnowledgeURLOption) {
		remove := true
		o.AutoRemove = &remove
	}
}

// Minimum frequency (in days) at which the document is refreshed. Defaults to 7.
func WithKnowledgeMinFrequencyDays(days int) CreateKnowledgeURLOption {
	return func(o *createKnowledgeURLOption) {
		o.MinimumFrequencyDays = &days
	}
}

type createKnowledgeURLRequest struct {
	URL                  string `json:"url"`
	Name                 string `json:"name,omitempty"`
	ParentFolderID       string `json:"parent_folder_id,omitempty"`
	EnableAutoSync       *bool  `json:"enable_auto_sync,omitempty"`
	AutoRemove           *bool  `json:"auto_remove,omitempty"`
	MinimumFrequencyDays *int   `json:"minimum_frequency_days,omitempty"`
}

func (o *createKnowledgeURLOption) body(pageURL string) createKnowledgeURLRequest {
	return createKnowledgeURLRequest{
		URL:                  pageURL,
		Name:                 o.Name,
		ParentFolderID:       o.ParentFolderID,
		EnableAutoSync:       o.EnableAutoSync,
		AutoRemove:           o.AutoRemove,
		MinimumFrequencyDays: o.MinimumFrequencyDays,
	}
}

type createKnowledgeFileOption struct {
	knowledgeNameFolderOption
	Filename string
}

type CreateKnowledgeFileOption func(*createKnowledgeFileOption)

// A custom, human-readable name for the document.
func WithKnowledgeFileName(name string) CreateKnowledgeFileOption {
	return func(o *createKnowledgeFileOption) {
		withKnowledgeName(name, &o.knowledgeNameFolderOption)
	}
}

// If set, the created document is placed inside the given folder.
func WithKnowledgeFileParentFolderID(id string) CreateKnowledgeFileOption {
	return func(o *createKnowledgeFileOption) {
		withKnowledgeParentFolderID(id, &o.knowledgeNameFolderOption)
	}
}

// Filename sent with the multipart file part. Defaults to the os.File name when available, otherwise "document".
func WithKnowledgeFilename(filename string) CreateKnowledgeFileOption {
	return func(o *createKnowledgeFileOption) {
		o.Filename = filename
	}
}

func (o *createKnowledgeFileOption) body(file io.Reader, filename string) ([]byte, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", err
	}
	n, err := io.Copy(fw, file)
	if err != nil {
		return nil, "", err
	}
	if n == 0 {
		return nil, "", fmt.Errorf("file is empty")
	}
	if o.Name != "" {
		if err := w.WriteField("name", o.Name); err != nil {
			return nil, "", err
		}
	}
	if o.ParentFolderID != "" {
		if err := w.WriteField("parent_folder_id", o.ParentFolderID); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), w.FormDataContentType(), nil
}

type getKnowledgeOption struct {
	AgentID string
}

type GetKnowledgeOption func(*getKnowledgeOption)

// Scope the document lookup to this agent.
func WithKnowledgeAgentID(id string) GetKnowledgeOption {
	return func(o *getKnowledgeOption) {
		o.AgentID = id
	}
}

func (o *getKnowledgeOption) query() url.Values {
	q := url.Values{}
	if o.AgentID != "" {
		q.Set("agent_id", o.AgentID)
	}
	return q
}

type updateKnowledgeOption struct {
	Name    *string
	Content *string
}

type UpdateKnowledgeOption func(*updateKnowledgeOption)

// A custom, human-readable name for the document.
func WithUpdateKnowledgeName(name string) UpdateKnowledgeOption {
	return func(o *updateKnowledgeOption) {
		o.Name = &name
	}
}

// Updated content for the document. Supported for text documents, URL documents with auto-sync disabled, and file documents that are not synced from an external source.
func WithUpdateKnowledgeContent(content string) UpdateKnowledgeOption {
	return func(o *updateKnowledgeOption) {
		o.Content = &content
	}
}

type updateKnowledgeRequest struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
}

func (o *updateKnowledgeOption) body() updateKnowledgeRequest {
	return updateKnowledgeRequest{
		Name:    o.Name,
		Content: o.Content,
	}
}

type deleteKnowledgeOption struct {
	Force *bool
}

type DeleteKnowledgeOption func(*deleteKnowledgeOption)

// If set to true, the document is deleted even if agents use it, and it is removed from those agents.
func WithForceDeleteKnowledge() DeleteKnowledgeOption {
	return func(o *deleteKnowledgeOption) {
		force := true
		o.Force = &force
	}
}

func (o *deleteKnowledgeOption) query() url.Values {
	q := url.Values{}
	if o.Force != nil {
		q.Set("force", strconv.FormatBool(*o.Force))
	}
	return q
}

type knowledgeDependentOption struct {
	Cursor   string
	PageSize int
}

type KnowledgeDependentOption func(*knowledgeDependentOption)

// Used for fetching the next page of dependent agents. Cursor is returned as NextCursor.
func WithKnowledgeDependentCursor(cursor string) KnowledgeDependentOption {
	return func(o *knowledgeDependentOption) {
		o.Cursor = cursor
	}
}

// How many dependent agents to return at maximum. Can not exceed 100, defaults to 30.
func WithKnowledgeDependentPageSize(size int) KnowledgeDependentOption {
	return func(o *knowledgeDependentOption) {
		o.PageSize = size
	}
}

func (o *knowledgeDependentOption) query() url.Values {
	q := url.Values{}
	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	return q
}

type getKnowledgeChunkOption struct {
	EmbeddingModel string
}

type GetKnowledgeChunkOption func(*getKnowledgeChunkOption)

// Embedding model used to look up the chunk.
func WithKnowledgeEmbeddingModel(model EmbeddingModel) GetKnowledgeChunkOption {
	return func(o *getKnowledgeChunkOption) {
		o.EmbeddingModel = model.String()
	}
}

func (o *getKnowledgeChunkOption) query() url.Values {
	q := url.Values{}
	if o.EmbeddingModel != "" {
		q.Set("embedding_model", o.EmbeddingModel)
	}
	return q
}

type ragIndexRequest struct {
	Model string `json:"model"`
}
