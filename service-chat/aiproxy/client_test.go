package aiproxy

import "testing"

func TestGetConversationTitleFromLLM(t *testing.T) {
	// 构造要发送给python端的消息
	send := ConversationTitleRequest{
		Content:        "今天天气怎么样呢?",
		ExcludedTitles: "天气查询、今日天气如何？、今天的天气情况如何？",
	}

	// 发送
	resp, err := GetConversationTitleFromLLM(&send)
	if err != nil {
		t.Errorf("GetConversationTitleFromLLM returned error:%v", err)
	}

	// 打印信息
	t.Log(resp.Content)
}

func TestGetMessageFromLLM(t *testing.T) {
	// 构建要发送给python端的消息
	send := ConversationMessageRequest{
		Content:          "假如给我三天光明这本书的主要内容",
		ConversationMode: "primary_assistant",
		ConversationID:   "1004",
	}

	// 发送
	resp, err := GetMessageFromLLM(&send)
	if err != nil {
		t.Errorf("GetMessageFromLLM returned error:%v", err)
	}

	// 打印信息
	t.Log(resp.Content)
	t.Log(resp.TokenUsage)
}
