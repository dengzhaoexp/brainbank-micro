package chatdtos

type ConversationRequest struct {
	ConversionID     string           `json:"conversation_id"`
	ConversationMode ConversationMode `json:"conversation_mode" binding:"required"`
	Message          Message          `json:"message" binding:"required"`
	Model            string           `json:"model" binding:"required"`
	Language         string           `json:"language" binding:"required"`
}

type ConversationMode struct {
	Kind string `json:"kind" binding:"required"`
}

type Message struct {
	Author  Author  `json:"author" binding:"required"`
	Content Content `json:"content" binding:"required"`
}

type Author struct {
	Role string `json:"role" binding:"required"`
}

type Content struct {
	ContentType string   `json:"content_json" binding:"required"`
	Parts       []string `json:"parts" binding:"required"`
}

type TitleResponse struct {
	Type           string `json:"type"`
	Title          string `json:"title"`
	ConversationID string `json:"conversation_id"`
}

type ConversationResponse struct {
	Message        Message `json:"message"`
	ConversationID string  `json:"conversation_id"`
}

type ConversationsResponse struct {
	Items                   []*ConversationItem `json:"items"`
	Total                   int                 `json:"total"`
	Limit                   int                 `json:"limit"`
	Offset                  int                 `json:"offset"`
	HasMissingConversations bool                `json:"has_missing_conversations"`
}

type ConversationItem struct {
	ConversationID         string `json:"id"`
	Title                  string `json:"title"`
	CreateTime             string `json:"create_time"`
	UpdateTime             string `json:"update_time"`
	Mapping                string `json:"mapping"`
	CurrentNode            string `json:"current_node"`
	ConversationTemplateId string `json:"conversation_template_id"`
	GizmoId                string `json:"gizmo_id"`
	IsArchived             bool   `json:"is_archived"`
	WorkspaceId            string `json:"workspace_id"`
}

type GetConversationResponse struct {
	ConversationID string  `json:"conversation_id"`
	Title          string  `json:"title"`
	CreateTime     string  `json:"create_time"`
	UpdateTime     string  `json:"update_time"`
	Mapping        Mapping `json:"mapping"`
}

type Mapping struct {
	Mapping map[string]MessageStruct `json:"mapping"`
}

type MessageStruct struct {
	ID         string  `json:"id"`
	CreateTime string  `json:"create_time"`
	Content    Content `json:"content"`
	Type       string  `json:"Type"`
}
