package aiproxy

type ConversationTitleRequest struct {
	Content        string `json:"content"`
	ExcludedTitles string `json:"excluded_titles"`
}

type ConversationTitleResponse struct {
	Content string `json:"content"`
}

type ConversationMessageRequest struct {
	Content          string `json:"content"`
	ConversationMode string `json:"conversation_mode"`
	ConversationID   string `json:"conversation_id"`
}

type ConversationMessageResponse struct {
	Content    []string `json:"content"`
	TokenUsage int64    `json:"token_usage"`
}
