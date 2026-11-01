package eleven

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestKnowledgeEnums(t *testing.T) {
	if KnowledgeTypeText.String() != "text" {
		t.Fatalf("type = %s", KnowledgeTypeText)
	}
	if KnowledgeSortUpdatedAt.String() != "updated_at" {
		t.Fatalf("sort = %s", KnowledgeSortUpdatedAt)
	}
	if KnowledgeUsagePrompt.String() != "prompt" {
		t.Fatalf("usage = %s", KnowledgeUsagePrompt)
	}
	if EmbeddingModelE5Mistral7BInstruct.String() != "e5_mistral_7b_instruct" {
		t.Fatalf("model = %s", EmbeddingModelE5Mistral7BInstruct)
	}
	if RAGIndexSucceeded.String() != "succeeded" {
		t.Fatalf("status = %s", RAGIndexSucceeded)
	}
	if RAGIndexCannotIndexFolder.String() != "cannot_index_folder" {
		t.Fatalf("status = %s", RAGIndexCannotIndexFolder)
	}
}

func TestListKnowledgeQuery(t *testing.T) {
	var opt listKnowledgeOption
	if q := opt.query().Encode(); q != "" {
		t.Fatalf("default query = %q", q)
	}

	WithKnowledgePageSize(20)(&opt)
	WithKnowledgeSearch("hours")(&opt)
	WithKnowledgeCreatedByUserID("@me")(&opt)
	WithKnowledgeTypes(KnowledgeTypeText, KnowledgeTypeURL)(&opt)
	WithKnowledgeParentFolderID("fld_1")(&opt)
	WithKnowledgeAncestorFolderID("fld_root")(&opt)
	WithKnowledgeFoldersFirst()(&opt)
	WithKnowledgeSort(KnowledgeSortName)(&opt)
	WithKnowledgeSortDirection(SortAsc)(&opt)
	WithKnowledgeCursor("cur_1")(&opt)

	q := opt.query()
	if q.Get("page_size") != "20" {
		t.Fatalf("page_size = %q", q.Get("page_size"))
	}
	if q.Get("search") != "hours" {
		t.Fatalf("search = %q", q.Get("search"))
	}
	if q.Get("created_by_user_id") != "@me" {
		t.Fatalf("created_by_user_id = %q", q.Get("created_by_user_id"))
	}
	if got := q["types"]; len(got) != 2 || got[0] != "text" || got[1] != "url" {
		t.Fatalf("types = %v", got)
	}
	if q.Get("parent_folder_id") != "fld_1" {
		t.Fatalf("parent_folder_id = %q", q.Get("parent_folder_id"))
	}
	if q.Get("ancestor_folder_id") != "fld_root" {
		t.Fatalf("ancestor_folder_id = %q", q.Get("ancestor_folder_id"))
	}
	if q.Get("folders_first") != "true" {
		t.Fatalf("folders_first = %q", q.Get("folders_first"))
	}
	if q.Get("sort_by") != "name" {
		t.Fatalf("sort_by = %q", q.Get("sort_by"))
	}
	if q.Get("sort_direction") != "asc" {
		t.Fatalf("sort_direction = %q", q.Get("sort_direction"))
	}
	if q.Get("cursor") != "cur_1" {
		t.Fatalf("cursor = %q", q.Get("cursor"))
	}
}

func TestCreateKnowledgeTextBody(t *testing.T) {
	var opt createKnowledgeTextOption
	got, err := json.Marshal(opt.body("Office hours are 9-17 UTC."))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"text":"Office hours are 9-17 UTC."}` {
		t.Fatalf("minimal body = %s", got)
	}

	WithKnowledgeName("hours")(&opt)
	WithKnowledgeTextParentFolderID("fld_1")(&opt)
	got, err = json.Marshal(opt.body("Office hours are 9-17 UTC."))
	if err != nil {
		t.Fatal(err)
	}
	var payload createKnowledgeTextRequest
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Name != "hours" || payload.ParentFolderID != "fld_1" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestCreateKnowledgeURLBody(t *testing.T) {
	var opt createKnowledgeURLOption
	WithKnowledgeURLName("docs")(&opt)
	WithKnowledgeURLParentFolderID("fld_1")(&opt)
	WithKnowledgeAutoSync()(&opt)
	WithKnowledgeAutoRemove()(&opt)
	WithKnowledgeMinFrequencyDays(3)(&opt)

	got, err := json.Marshal(opt.body("https://example.com/docs"))
	if err != nil {
		t.Fatal(err)
	}
	var payload createKnowledgeURLRequest
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.URL != "https://example.com/docs" || payload.Name != "docs" {
		t.Fatalf("payload = %+v", payload)
	}
	if payload.EnableAutoSync == nil || !*payload.EnableAutoSync {
		t.Fatalf("auto_sync = %v", payload.EnableAutoSync)
	}
	if payload.MinimumFrequencyDays == nil || *payload.MinimumFrequencyDays != 3 {
		t.Fatalf("frequency = %v", payload.MinimumFrequencyDays)
	}
}

func TestCreateKnowledgeFileBody(t *testing.T) {
	var opt createKnowledgeFileOption
	WithKnowledgeFileName("faq")(&opt)
	WithKnowledgeFileParentFolderID("fld_1")(&opt)

	body, contentType, err := opt.body(strings.NewReader("faq text"), "faq.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data; boundary=") {
		t.Fatalf("content type = %q", contentType)
	}

	req, err := http.NewRequest(http.MethodPost, "http://example", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)
	if err := req.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	if req.FormValue("name") != "faq" || req.FormValue("parent_folder_id") != "fld_1" {
		t.Fatalf("form = %v", req.MultipartForm.Value)
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if header.Filename != "faq.txt" {
		t.Fatalf("filename = %q", header.Filename)
	}
	got, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "faq text" {
		t.Fatalf("file = %q", got)
	}
}

func TestUpdateKnowledgeBody(t *testing.T) {
	var opt updateKnowledgeOption
	got, err := json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{}` {
		t.Fatalf("empty patch = %s", got)
	}

	WithUpdateKnowledgeName("hours v2")(&opt)
	WithUpdateKnowledgeContent("Office hours are 10-18 UTC.")(&opt)
	got, err = json.Marshal(opt.body())
	if err != nil {
		t.Fatal(err)
	}
	var payload updateKnowledgeRequest
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Name == nil || *payload.Name != "hours v2" {
		t.Fatalf("name = %v", payload.Name)
	}
	if payload.Content == nil || *payload.Content != "Office hours are 10-18 UTC." {
		t.Fatalf("content = %v", payload.Content)
	}
}

func TestKnowledgeValidation(t *testing.T) {
	client := NewClient("test")
	if _, err := client.CreateKnowledgeBaseText(""); err == nil || err.Error() != "text is empty" {
		t.Fatalf("text err = %v", err)
	}
	if _, err := client.CreateKnowledgeBaseURL(""); err == nil || err.Error() != "url is empty" {
		t.Fatalf("url err = %v", err)
	}
	if _, err := client.CreateKnowledgeBaseFile(nil); err == nil || err.Error() != "file is nil" {
		t.Fatalf("file nil err = %v", err)
	}
	if _, err := client.CreateKnowledgeBaseFile(strings.NewReader("")); err == nil || err.Error() != "file is empty" {
		t.Fatalf("file empty err = %v", err)
	}
	if _, err := client.GetKnowledgeBase(""); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("get err = %v", err)
	}
	if _, err := client.UpdateKnowledgeBase(""); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("update err = %v", err)
	}
	if err := client.DeleteKnowledgeBase(""); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("delete err = %v", err)
	}
	if _, err := client.GetKnowledgeBaseContent(""); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("content err = %v", err)
	}
	if _, err := client.GetKnowledgeBaseDependentAgents(""); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("dependent err = %v", err)
	}
	if _, err := client.ComputeRAGIndex("", EmbeddingModelE5Mistral7BInstruct); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("rag id err = %v", err)
	}
	if _, err := client.ComputeRAGIndex("doc_1", ""); err == nil || err.Error() != "model is empty" {
		t.Fatalf("rag model err = %v", err)
	}
	if _, err := client.GetKnowledgeBaseChunk("", "chunk_1"); err == nil || err.Error() != "documentationID is empty" {
		t.Fatalf("chunk doc err = %v", err)
	}
	if _, err := client.GetKnowledgeBaseChunk("doc_1", ""); err == nil || err.Error() != "chunkID is empty" {
		t.Fatalf("chunk id err = %v", err)
	}
}

func TestListKnowledgeBase(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		if r.URL.Query().Get("search") != "hours" || r.URL.Query().Get("page_size") != "20" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"documents":[{"type":"text","id":"doc_1","name":"hours","access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"metadata":{"created_at_unix_secs":1,"last_updated_at_unix_secs":2,"size_bytes":12},"supported_usages":["auto","prompt"]}],"has_more":true,"next_cursor":"cur_2"}`))
	})

	got, err := client.ListKnowledgeBase(WithKnowledgeSearch("hours"), WithKnowledgePageSize(20))
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasMore || got.NextCursor != "cur_2" || len(got.Documents) != 1 {
		t.Fatalf("list = %+v", got)
	}
	if got.Documents[0].ID != "doc_1" || got.Documents[0].Type != KnowledgeTypeText {
		t.Fatalf("document = %+v", got.Documents[0])
	}
}

func TestCreateKnowledgeBaseText(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/text" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		requireHeader(t, r, xiAPIKeyHeader, "test-key")
		var payload createKnowledgeTextRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Text != "Office hours are 9-17 UTC." || payload.Name != "hours" {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"id":"doc_1","name":"hours"}`))
	})

	got, err := client.CreateKnowledgeBaseText("Office hours are 9-17 UTC.", WithKnowledgeName("hours"))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "doc_1" || got.Name != "hours" {
		t.Fatalf("result = %+v", got)
	}
}

func TestCreateKnowledgeBaseURL(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/knowledge-base/url" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload createKnowledgeURLRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.URL != "https://example.com/docs" || payload.Name != "docs" {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"id":"doc_2","name":"docs","folder_path":[{"id":"fld_root"}]}`))
	})

	got, err := client.CreateKnowledgeBaseURL("https://example.com/docs", WithKnowledgeURLName("docs"))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "doc_2" || len(got.FolderPath) != 1 {
		t.Fatalf("result = %+v", got)
	}
}

func TestCreateKnowledgeBaseFile(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/file" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("content-type = %q", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("name") != "faq" {
			t.Fatalf("name = %q", r.FormValue("name"))
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		if header.Filename != "faq.txt" {
			t.Fatalf("filename = %q", header.Filename)
		}
		got, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "faq text" {
			t.Fatalf("file = %q", got)
		}
		w.Write([]byte(`{"id":"doc_3","name":"faq"}`))
	})

	got, err := client.CreateKnowledgeBaseFile(
		strings.NewReader("faq text"),
		WithKnowledgeFileName("faq"),
		WithKnowledgeFilename("faq.txt"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "doc_3" {
		t.Fatalf("result = %+v", got)
	}
}

func TestGetKnowledgeBase(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("agent_id") != "agent_1" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"type":"text","id":"doc_1","name":"hours","extracted_inner_html":"<p>hours</p>","access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"metadata":{"created_at_unix_secs":1,"last_updated_at_unix_secs":2,"size_bytes":12},"supported_usages":["auto"]}`))
	})

	got, err := client.GetKnowledgeBase("doc_1", WithKnowledgeAgentID("agent_1"))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "doc_1" || got.ExtractedHTML != "<p>hours</p>" {
		t.Fatalf("document = %+v", got)
	}
}

func TestUpdateKnowledgeBase(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload updateKnowledgeRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Name == nil || *payload.Name != "hours v2" {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"type":"text","id":"doc_1","name":"hours v2","access_info":{"is_creator":true,"creator_name":"Ada","creator_email":"ada@example.com","role":"admin"},"metadata":{"created_at_unix_secs":1,"last_updated_at_unix_secs":3,"size_bytes":14},"supported_usages":["auto"]}`))
	})

	got, err := client.UpdateKnowledgeBase("doc_1", WithUpdateKnowledgeName("hours v2"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "hours v2" {
		t.Fatalf("document = %+v", got)
	}
}

func TestDeleteKnowledgeBase(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("force") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{}`))
	})
	if err := client.DeleteKnowledgeBase("doc_1", WithForceDeleteKnowledge()); err != nil {
		t.Fatal(err)
	}
}

func TestGetKnowledgeBaseContent(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1/content" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`"Office hours are 9-17 UTC."`))
	})

	got, err := client.GetKnowledgeBaseContent("doc_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Content != "Office hours are 9-17 UTC." {
		t.Fatalf("content = %+v", got)
	}
}

func TestDecodeKnowledgeContent(t *testing.T) {
	if got := decodeKnowledgeContent([]byte(`"hello"`)); got.Content != "hello" {
		t.Fatalf("string = %+v", got)
	}
	if got := decodeKnowledgeContent([]byte(`{"content":"from object"}`)); got.Content != "from object" {
		t.Fatalf("object = %+v", got)
	}
	if got := decodeKnowledgeContent([]byte(`plain text`)); got.Content != "plain text" {
		t.Fatalf("raw = %+v", got)
	}
}

func TestGetKnowledgeBaseDependentAgents(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1/dependent-agents" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("cursor") != "cur_1" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"agents":[{"type":"unknown","id":"agent_x"}],"has_more":false}`))
	})

	got, err := client.GetKnowledgeBaseDependentAgents("doc_1", WithKnowledgeDependentCursor("cur_1"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Agents) != 1 || got.Agents[0].Type != DependentAgentUnknown {
		t.Fatalf("agents = %+v", got.Agents)
	}
}

func TestComputeRAGIndex(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1/rag-index" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var payload ragIndexRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Model != "e5_mistral_7b_instruct" {
			t.Fatalf("payload = %+v", payload)
		}
		w.Write([]byte(`{"id":"idx_1","model":"e5_mistral_7b_instruct","status":"succeeded","progress_percentage":100,"document_model_index_usage":{"used_bytes":128}}`))
	})

	got, err := client.ComputeRAGIndex("doc_1", EmbeddingModelE5Mistral7BInstruct)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "idx_1" || got.Status != RAGIndexSucceeded || got.DocumentModelIndexUsage.UsedBytes != 128 {
		t.Fatalf("index = %+v", got)
	}
}

func TestGetKnowledgeBaseChunk(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/knowledge-base/doc_1/chunk/chunk_1" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("embedding_model") != "e5_mistral_7b_instruct" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"id":"chunk_1","name":"hours","content":"Office hours are 9-17 UTC."}`))
	})

	got, err := client.GetKnowledgeBaseChunk("doc_1", "chunk_1", WithKnowledgeEmbeddingModel(EmbeddingModelE5Mistral7BInstruct))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "chunk_1" || got.Content != "Office hours are 9-17 UTC." {
		t.Fatalf("chunk = %+v", got)
	}
}

func TestKnowledgeError(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	})
	if _, err := client.GetKnowledgeBase("missing"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := client.CreateKnowledgeBaseText("hours"); err == nil {
		t.Fatal("expected error")
	}
}

func TestKnowledgeDocumentsAll(t *testing.T) {
	docs := KnowledgeDocuments{
		{ID: "a", Type: KnowledgeTypeText},
		{ID: "b", Type: KnowledgeTypeURL},
	}

	var ids []string
	for document := range docs.All() {
		ids = append(ids, document.ID)
	}
	if len(ids) != 2 {
		t.Fatalf("ids = %v", ids)
	}

	var urls []string
	for document := range docs.AllWith(func(d *KnowledgeDocument) bool {
		return d.Type == KnowledgeTypeURL
	}) {
		urls = append(urls, document.ID)
	}
	if len(urls) != 1 || urls[0] != "b" {
		t.Fatalf("urls = %v", urls)
	}
}

func TestAgentKnowledgeBaseLocatorJSON(t *testing.T) {
	prompt := AgentPrompt{
		ToolIDs: []string{"tool_1"},
		KnowledgeBase: []AgentKnowledgeBaseLocator{{
			Type:      KnowledgeTypeText,
			Name:      "hours",
			ID:        "doc_1",
			UsageMode: KnowledgeUsageAuto,
		}},
		RAG: json.RawMessage(`{"enabled":true}`),
	}
	got, err := json.Marshal(prompt)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(got, &payload); err != nil {
		t.Fatal(err)
	}
	kb := payload["knowledge_base"].([]any)
	if len(kb) != 1 {
		t.Fatalf("knowledge_base = %v", payload["knowledge_base"])
	}
	item := kb[0].(map[string]any)
	if item["type"] != "text" || item["id"] != "doc_1" || item["usage_mode"] != "auto" {
		t.Fatalf("item = %v", item)
	}
}

func TestLiveKnowledgeBase(t *testing.T) {
	client := liveClient(t)
	created, err := client.CreateKnowledgeBaseText(
		"Office hours for the go-eleven SDK live test are 9-17 UTC.",
		WithKnowledgeName("go-eleven-knowledge-test"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" {
		t.Fatal("empty documentation id")
	}
	t.Cleanup(func() {
		_ = client.DeleteKnowledgeBase(created.ID, WithForceDeleteKnowledge())
	})

	got, err := client.GetKnowledgeBase(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID || got.Name != "go-eleven-knowledge-test" {
		t.Fatalf("get = %+v", got)
	}

	list, err := client.ListKnowledgeBase(WithKnowledgeSearch("go-eleven-knowledge-test"), WithKnowledgePageSize(30))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for document := range list.Documents.All() {
		if document.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created document missing from list: %+v", list)
	}

	index, err := client.ComputeRAGIndex(created.ID, EmbeddingModelE5Mistral7BInstruct)
	if err != nil {
		t.Fatal(err)
	}
	if index.Status == "" {
		t.Fatalf("index = %+v", index)
	}

	if err := client.DeleteKnowledgeBase(created.ID, WithForceDeleteKnowledge()); err != nil {
		t.Fatal(err)
	}
}
