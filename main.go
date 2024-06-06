package main

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	MessageID         string    `json:"message_id"`
	SessionID         string    `json:"session_id"`
	SenderID          string    `json:"sender"`
	SendAt            time.Time `json:"send_at"`
	MessageType       string    `json:"message_type"` // text、pic
	Content           string    `json:"content"`
	PreviousMessageId string    `json:"previous_message_id"`
	ParentMessageId   string    `json:"parent_message_id"`
}

type SessionPromptTemplate struct {
	TemplateID           string            `json:"template_id"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	PromptText           string            `json:"prompt_text"`
	Variables            []string          `json:"variables"`
	VariableDescriptions map[string]string `json:"variable_descriptions"`
	Tags                 []string          `json:"tags"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	CreatedBy            string            `json:"created_by"`
	UpdatedBy            string            `json:"updated_by"`
	Version              float32           `json:"version"`
	Status               string            `json:"status"`
}

func main() {
	template := &SessionPromptTemplate{
		TemplateID:           "c748af67-fcdd-856a-de7e-c37983fede6c",
		Name:                 "session_context",
		Description:          "Prompt template for generating session's context",
		PromptText:           "根据用户提出的问题\"{message}\",生成一段富有见解的{language}上下文信息,作为本次会话的背景知识。上下文应当包含以下几个方面:\n1. 对问题主题的概括性介绍\n2. 问题主题的重要性和应用场景 \n3. 相关的基础概念和原理\n4. 值得了解的有趣事实或见解\n\n会话上下文: {context}",
		Variables:            []string{"language", "message", "context"},
		VariableDescriptions: map[string]string{"language": "语言", "message": "用户发送的消息", "context": "聊天的背景"},
		Tags:                 []string{"context-generation", "prompt-template", "summarization", "chatdtos-service"},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CreatedBy:            "dunzane",
		UpdatedBy:            "dunzane",
		Version:              0.0,
		Status:               "sketch",
	}

	jsonData, _ := json.MarshalIndent(template, "", "    ")
	jsonStr := string(jsonData)
	print(jsonStr)
}
