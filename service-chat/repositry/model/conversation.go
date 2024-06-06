package model

import (
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ConversationId string `gorm:"primarykey"`
	Title          string
	Status         string // Created、Activate、Idle、Expired、Abandoned、Renewed、Terminated、Invalidated、Migrated
	CreatedBy      string
	MessageCount   int
	Language       string
	TokensConsumed int64
}

//type PromptTemplate struct {
//	TemplateId           string            `json:"template_id"`
//	Name                 string            `json:"name"`
//	Description          string            `json:"description"`
//	PromptText           string            `json:"prompt_text"`
//	Variables            []string          `json:"variables"`
//	VariableDescriptions map[string]string `json:"variable_descriptions"`
//	Tags                 []string          `json:"tags"`
//	CreatedAt            time.Time         `json:"created_at"`
//	UpdatedAt            time.Time         `json:"updated_at"`
//	CreatedBy            string            `json:"created_by"`
//	UpdatedBy            string            `json:"updated_by"`
//	Version              float32           `json:"version"`
//	Status               string            `json:"status"`
//}
