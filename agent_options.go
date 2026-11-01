package eleven

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type listAgentsOption struct {
	PageSize        int
	Search          string
	Archived        *bool
	CreatedByUserID string
	Tags            []string
	Sort            string
	SortDirection   string
	Cursor          string
}

type ListAgentsOption func(*listAgentsOption)

// How many Agents to return at maximum. Can not exceed 100, defaults to 30.
func WithAgentPageSize(size int) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.PageSize = size
	}
}

// Search by agents name.
func WithAgentSearch(search string) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.Search = search
	}
}

// Filter agents by archived status. Sends archived=true.
func WithAgentArchived() ListAgentsOption {
	return func(o *listAgentsOption) {
		archived := true
		o.Archived = &archived
	}
}

// Filter agents by creator user ID. When set, only agents created by this user are returned. Use "@me" to refer to the authenticated user.
func WithCreatedByUserID(id string) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.CreatedByUserID = id
	}
}

// Filter agents by tag. Repeat the parameter to match any of several tags.
func WithAgentTags(tags ...string) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.Tags = tags
	}
}

// The field to sort the results by. name, created_at, updated_at, or call_count_7d.
func WithAgentSort(sort AgentSort) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.Sort = sort.String()
	}
}

// The direction to sort the results. asc or desc.
func WithAgentSortDirection(direction SortDirection) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.SortDirection = direction.String()
	}
}

// Used for fetching the next page. Cursor is returned in the list response as NextCursor.
func WithAgentCursor(cursor string) ListAgentsOption {
	return func(o *listAgentsOption) {
		o.Cursor = cursor
	}
}

func (o *listAgentsOption) query() url.Values {
	q := url.Values{}
	if o.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.Search != "" {
		q.Set("search", o.Search)
	}
	if o.Archived != nil {
		q.Set("archived", strconv.FormatBool(*o.Archived))
	}
	if o.CreatedByUserID != "" {
		q.Set("created_by_user_id", o.CreatedByUserID)
	}
	for _, tag := range o.Tags {
		q.Add("tags", tag)
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

type getAgentOption struct {
	VersionID                         string
	BranchID                          string
	IncludeProcedureGeneratedWorkflow *bool
}

type GetAgentOption func(*getAgentOption)

// The ID of the agent version to use.
func WithAgentVersionID(id string) GetAgentOption {
	return func(o *getAgentOption) {
		o.VersionID = id
	}
}

// The ID of the branch to use.
func WithAgentBranchID(id string) GetAgentOption {
	return func(o *getAgentOption) {
		o.BranchID = id
	}
}

// When false, the workflow contains only authored elements, without nodes and edges generated from structured procedures.
func WithoutProcedureGeneratedWorkflow() GetAgentOption {
	return func(o *getAgentOption) {
		include := false
		o.IncludeProcedureGeneratedWorkflow = &include
	}
}

func (o *getAgentOption) query() url.Values {
	q := url.Values{}
	if o.VersionID != "" {
		q.Set("version_id", o.VersionID)
	}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.IncludeProcedureGeneratedWorkflow != nil {
		q.Set("include_procedure_generated_workflow_elements", strconv.FormatBool(*o.IncludeProcedureGeneratedWorkflow))
	}
	return q
}

type createAgentOption struct {
	Name             string
	Tags             []string
	PlatformSettings json.RawMessage
	Workflow         json.RawMessage
}

type CreateAgentOption func(*createAgentOption)

// A name to make the agent easier to find.
func WithAgentName(name string) CreateAgentOption {
	return func(o *createAgentOption) {
		o.Name = name
	}
}

// Tags to help classify and filter the agent.
func WithCreateAgentTags(tags ...string) CreateAgentOption {
	return func(o *createAgentOption) {
		o.Tags = tags
	}
}

// Platform settings for the agent are all settings that aren't related to the conversation orchestration and content. Pass through JSON.
func WithPlatformSettings(settings json.RawMessage) CreateAgentOption {
	return func(o *createAgentOption) {
		o.PlatformSettings = settings
	}
}

// Workflow for the agent. Pass through JSON for the node and edge graph.
func WithAgentWorkflow(workflow json.RawMessage) CreateAgentOption {
	return func(o *createAgentOption) {
		o.Workflow = workflow
	}
}

type createAgentRequest struct {
	ConversationConfig ConversationalConfig `json:"conversation_config"`
	PlatformSettings   json.RawMessage      `json:"platform_settings,omitempty"`
	Workflow           json.RawMessage      `json:"workflow,omitempty"`
	Name               string               `json:"name,omitempty"`
	Tags               []string             `json:"tags,omitempty"`
}

func (o *createAgentOption) body(config ConversationalConfig) createAgentRequest {
	return createAgentRequest{
		ConversationConfig: config,
		PlatformSettings:   o.PlatformSettings,
		Workflow:           o.Workflow,
		Name:               o.Name,
		Tags:               o.Tags,
	}
}

type updateAgentOption struct {
	Name                              *string
	Tags                              []string
	ConversationConfig                *ConversationalConfig
	PlatformSettings                  json.RawMessage
	Workflow                          json.RawMessage
	VersionDescription                *string
	BranchID                          string
	IncludeProcedureGeneratedWorkflow *bool
}

type UpdateAgentOption func(*updateAgentOption)

// A name to make the agent easier to find.
func WithUpdateAgentName(name string) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.Name = &name
	}
}

// Replace agent tags.
func WithUpdateAgentTags(tags ...string) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.Tags = tags
	}
}

// Conversation configuration for an agent.
func WithUpdateConversationConfig(config ConversationalConfig) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.ConversationConfig = &config
	}
}

// Platform settings for the agent. Pass through JSON.
func WithUpdatePlatformSettings(settings json.RawMessage) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.PlatformSettings = settings
	}
}

// Workflow for the agent. Pass through JSON.
func WithUpdateAgentWorkflow(workflow json.RawMessage) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.Workflow = workflow
	}
}

// Description for this version when publishing changes.
func WithVersionDescription(description string) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.VersionDescription = &description
	}
}

// The ID of the branch to use.
func WithUpdateAgentBranchID(id string) UpdateAgentOption {
	return func(o *updateAgentOption) {
		o.BranchID = id
	}
}

// When false, the workflow contains only authored elements, without nodes and edges generated from structured procedures.
func WithoutUpdateProcedureGeneratedWorkflow() UpdateAgentOption {
	return func(o *updateAgentOption) {
		include := false
		o.IncludeProcedureGeneratedWorkflow = &include
	}
}

type updateAgentRequest struct {
	ConversationConfig *ConversationalConfig `json:"conversation_config,omitempty"`
	PlatformSettings   json.RawMessage       `json:"platform_settings,omitempty"`
	Workflow           json.RawMessage       `json:"workflow,omitempty"`
	Name               *string               `json:"name,omitempty"`
	Tags               []string              `json:"tags,omitempty"`
	VersionDescription *string               `json:"version_description,omitempty"`
}

func (o *updateAgentOption) query() url.Values {
	q := url.Values{}
	if o.BranchID != "" {
		q.Set("branch_id", o.BranchID)
	}
	if o.IncludeProcedureGeneratedWorkflow != nil {
		q.Set("include_procedure_generated_workflow_elements", strconv.FormatBool(*o.IncludeProcedureGeneratedWorkflow))
	}
	return q
}

func (o *updateAgentOption) body() updateAgentRequest {
	return updateAgentRequest{
		ConversationConfig: o.ConversationConfig,
		PlatformSettings:   o.PlatformSettings,
		Workflow:           o.Workflow,
		Name:               o.Name,
		Tags:               o.Tags,
		VersionDescription: o.VersionDescription,
	}
}

type deleteAgentOption struct {
	Force *bool
}

type DeleteAgentOption func(*deleteAgentOption)

// If set to true, the agent is deleted even if other agents use it as a subagent, and it is removed from those agents and their branches.
func WithForceDeleteAgent() DeleteAgentOption {
	return func(o *deleteAgentOption) {
		force := true
		o.Force = &force
	}
}

func (o *deleteAgentOption) query() url.Values {
	q := url.Values{}
	if o.Force != nil {
		q.Set("force", strconv.FormatBool(*o.Force))
	}
	return q
}

type duplicateAgentOption struct {
	Name string
}

type DuplicateAgentOption func(*duplicateAgentOption)

// A name to make the duplicated agent easier to find.
func WithDuplicateAgentName(name string) DuplicateAgentOption {
	return func(o *duplicateAgentOption) {
		o.Name = name
	}
}

type duplicateAgentRequest struct {
	Name string `json:"name,omitempty"`
}

func (o *duplicateAgentOption) body() duplicateAgentRequest {
	return duplicateAgentRequest{Name: o.Name}
}
