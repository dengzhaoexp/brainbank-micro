package aiproxy

import (
	"encoding/json"
)

func GetConversationTitleFromLLM(req *ConversationTitleRequest) (*ConversationTitleResponse, error) {
	// 构建请求体
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 发送POST请求
	resp, err := postRequest("http://127.0.0.1:8000/conversation/title", requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应体
	var titleResponse ConversationTitleResponse
	if err = json.NewDecoder(resp.Body).Decode(&titleResponse); err != nil {
		return nil, err
	}

	return &titleResponse, nil
}

func GetMessageFromLLM(req *ConversationMessageRequest) (*ConversationMessageResponse, error) {
	// 构建请求体
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 发送POST请求
	resp, err := postRequest("http://127.0.0.1:8000/conversation/message", requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应体
	var messageResponse ConversationMessageResponse
	if err = json.NewDecoder(resp.Body).Decode(&messageResponse); err != nil {
		return nil, err
	}

	return &messageResponse, nil
}
