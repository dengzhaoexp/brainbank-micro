package model

import (
	"time"
)

type Neuron struct {
	NeuronID      string     `gorm:"primaryKey"` // 唯一标识符
	Name          string     // 神经元名称
	Description   string     // 描述
	OwnedBy       string     // 所以者（用户ID）
	CreatedAt     *time.Time `gorm:"default:null"` // 创建时间
	UpdatedAt     *time.Time `gorm:"default:null"` // 更新时间
	DocumentCount int        // 文档数量
	SessionCount  int        // 会话数量
	TokenUsage    int        // Token使用量
}

type NeuronConfig struct {
	ConfigID          string // 唯一标识符
	NeuronID          uint   // 指向 Neuron
	LLMType           string // 使用的大语言模型类型
	IndexingStrategy  string // 知识库索引策略
	RetrievalStrategy string // 知识检索策略
}

type LLMParam struct {
	ParamID  string `gorm:"autoIncrement"`
	ConfigID string
	Name     string
	Value    string
}

type Permission struct {
	PermissionID string `gorm:"autoIncrement"`
	NeuronID     uint
	UserID       string // 用户ID
	Role         string // 角色(所有者、读写、只读等)
}
