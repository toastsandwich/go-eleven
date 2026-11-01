package eleven

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type listToolsOption struct {
	Search          string
	PageSize        int
	CreatedByUserID string
	Types           []string
	Sort            string
	SortDirection   string
	UsedByAgentID   string
	Cursor          string
}

type ListToolsOption func(*listToolsOption)

// If specified, the endpoint returns only tools whose names start with this string.
func WithToolSearch(search string) ListToolsOption {
	return func(o *listToolsOption) {
		o.Search = search
	}
}

// How many tools to return at maximum. Can not exceed 100, defaults to 30.
func WithToolPageSize(size int) ListToolsOption {
	return func(o *listToolsOption) {
		o.PageSize = size
	}
}

// Filter tools by creator user ID. When set, only tools created by this user are returned. Use "@me" to refer to the authenticated user.
func WithToolCreatedByUserID(id string) ListToolsOption {
	return func(o *listToolsOption) {
		o.CreatedByUserID = id
	}
}

// If present, the endpoint returns only tools of the given types. webhook, client, or api_integration_webhook.
func WithToolTypes(types ...ToolType) ListToolsOption {
	return func(o *listToolsOption) {
		o.Types = make([]string, len(types))
		for i, t := range types {
			o.Types[i] = t.String()
		}
	}
}

// The field to sort the results by. name or created_at.
func WithToolSort(sort ToolSort) ListToolsOption {
	return func(o *listToolsOption) {
		o.Sort = sort.String()
	}
}

// The direction to sort the results. asc or desc.
func WithToolSortDirection(direction SortDirection) ListToolsOption {
	return func(o *listToolsOption) {
		o.SortDirection = direction.String()
	}
}

// Filter to tools referenced by this agent's main branch.
func WithToolUsedByAgentID(id string) ListToolsOption {
	return func(o *listToolsOption) {
		o.UsedByAgentID = id
	}
}

// Used for fetching the next page. Cursor is returned in the list response as NextCursor.
func WithToolCursor(cursor string) ListToolsOption {
	return func(o *listToolsOption) {
		o.Cursor = cursor
	}
}

func (o *listToolsOption) query() url.Values {
	q := url.Values{}
	if o.Search != "" {
		q.Set("search", o.Search)
	}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.CreatedByUserID != "" {
		q.Set("created_by_user_id", o.CreatedByUserID)
	}
	for _, t := range o.Types {
		q.Add("types", t)
	}
	if o.Sort != "" {
		q.Set("sort_by", o.Sort)
	}
	if o.SortDirection != "" {
		q.Set("sort_direction", o.SortDirection)
	}
	if o.UsedByAgentID != "" {
		q.Set("used_by_agent_id", o.UsedByAgentID)
	}
	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	return q
}

type getToolOption struct {
	Environment string
}

type GetToolOption func(*getToolOption)

// Environment whose values are used when the MCP server URL, headers, or auth connection reference environment variables. Defaults to production.
func WithToolEnvironment(environment string) GetToolOption {
	return func(o *getToolOption) {
		o.Environment = environment
	}
}

func (o *getToolOption) query() url.Values {
	q := url.Values{}
	if o.Environment != "" {
		q.Set("environment", o.Environment)
	}
	return q
}

type createToolOption struct {
	ResponseMocks json.RawMessage
}

type CreateToolOption func(*createToolOption)

// Mock responses with optional parameter conditions. Evaluated top-to-bottom; first match wins.
func WithToolResponseMocks(mocks json.RawMessage) CreateToolOption {
	return func(o *createToolOption) {
		o.ResponseMocks = mocks
	}
}

type createToolRequest struct {
	ToolConfig    ToolConfig      `json:"tool_config"`
	ResponseMocks json.RawMessage `json:"response_mocks,omitempty"`
}

func (o *createToolOption) body(config ToolConfig) createToolRequest {
	return createToolRequest{
		ToolConfig:    config,
		ResponseMocks: o.ResponseMocks,
	}
}

type UpdateToolOption = CreateToolOption

type deleteToolOption struct {
	Force *bool
}

type DeleteToolOption func(*deleteToolOption)

// If set to true, the tool is deleted even if agents use it, and it is removed from those agents and their branches.
func WithForceDeleteTool() DeleteToolOption {
	return func(o *deleteToolOption) {
		force := true
		o.Force = &force
	}
}

func (o *deleteToolOption) query() url.Values {
	q := url.Values{}
	if o.Force != nil {
		q.Set("force", strconv.FormatBool(*o.Force))
	}
	return q
}

type toolDependentOption struct {
	Cursor   string
	PageSize int
}

type ToolDependentOption func(*toolDependentOption)

// Used for fetching the next page of dependent agents. Cursor is returned as NextCursor.
func WithToolDependentCursor(cursor string) ToolDependentOption {
	return func(o *toolDependentOption) {
		o.Cursor = cursor
	}
}

// How many dependent agents to return at maximum. Can not exceed 100, defaults to 30.
func WithToolDependentPageSize(size int) ToolDependentOption {
	return func(o *toolDependentOption) {
		o.PageSize = size
	}
}

func (o *toolDependentOption) query() url.Values {
	q := url.Values{}
	if o.Cursor != "" {
		q.Set("cursor", o.Cursor)
	}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	return q
}
