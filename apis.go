package eleven

const (
	parentEndpoint = "https://api.elevenlabs.io/"
	xiAPIKeyHeader = "xi-api-key"
)
const (
	ttsEndpoint                   = "v1/text-to-speech"
	ttsStreamSegment              = "stream"
	speechToSpeechEndpoint        = "v1/voice-changer"
	speechToTextEndpoint          = "v1/speech-to-text"
	speechToTextTranscriptsPath   = "v1/speech-to-text/transcripts"
	listModelsEndpoint            = "v1/models"
	searchVoicesEndpoint          = "v2/voices"
	workspaceWebhooksEndpoint     = "v1/workspace/webhooks"
	agentsEndpoint                = "v1/convai/agents"
	createAgentEndpoint           = "v1/convai/agents/create"
	conversationsEndpoint         = "v1/convai/conversations"
	conversationSignedURLEndpoint = "v1/convai/conversation/get-signed-url"
	conversationTokenEndpoint     = "v1/convai/conversation/token"
	conversationWSPath            = "v1/convai/conversation"
	conversationWSSource          = "go-eleven"
	realtimeClientSecretsEndpoint = "v1/convai/realtime/client_secrets"
	toolsEndpoint                 = "v1/convai/tools"
	knowledgeBaseEndpoint         = "v1/convai/knowledge-base"
	knowledgeBaseTextEndpoint     = "v1/convai/knowledge-base/text"
	knowledgeBaseURLEndpoint      = "v1/convai/knowledge-base/url"
	knowledgeBaseFileEndpoint     = "v1/convai/knowledge-base/file"
	defaultSpeechToTextModel      = "scribe_v2"
)
