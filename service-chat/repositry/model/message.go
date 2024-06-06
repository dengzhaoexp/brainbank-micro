package model

import "time"

type Message struct {
	MessageID         string    `json:"messageId"`
	ConversationID    string    `json:"conversationId"`
	SenderID          string    `json:"senderId"`
	SenderName        string    `json:"senderName"`
	SentAt            time.Time `json:"sentAt"`
	MessageType       string    `json:"messageType"`
	Content           string    `json:"content"`
	Vector            []float64 `json:"vector,omitempty"`  // 可选的向量表示
	PreviousMessageID *string   `json:"previousMessageId"` // 对于开启会话的首条消息，该字段为null
	NextMessageID     *string   `json:"nextMessageId"`     // 后续消息的ID
}
