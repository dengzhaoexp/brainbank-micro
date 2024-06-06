package es

import (
	"chat/pkg/logger"
	"testing"
)

func TestGetMessage(t *testing.T) {
	// 加载基本信息
	logger.InitLog()
	InitElastic()

	_, err := GetMessage("chat-history", "7ded711d-a118-4548-9f6a-563928a1d8e0")
	if err != nil {
		return
	}

}
