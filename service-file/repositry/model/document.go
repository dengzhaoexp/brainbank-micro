package model

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string
	Type        string   // 例如 pdf, docx, txt 等
	Extension   string   // 文件扩展名
	Size        int64    // 文件大小(字节)
	ContentHash string   // 内容哈希值,用于检测重复文件
	Metadata    Metadata `gorm:"embedded"` // 文件元数据,如作者、标题等
	UploadedAt  string
	UploadedBy  string
	Status      string // 文件状态,如 Active, Deleted, Released 等
}

type Metadata struct {
	Title   string
	Author  string
	Subject string
	Created int64 // 这个字段暂且保留
}
