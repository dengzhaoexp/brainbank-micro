package es

import (
	"chat/pkg/logger"
	"encoding/json"
	"errors"
	"strings"
)

func GetMessageID(index, id string) (string, error) {
	// 设置查询体
	query := `{"query":{"match":{"session_id":"` + id + `"}},"sort":[{"created_at":{"order":"asc"}}]}`

	// 根据查询体查询出具体的内容
	search, err := _elastic.Search(
		_elastic.Search.WithIndex(index),
		_elastic.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return "", err
	}

	// 解析消息体
	var result map[string]interface{}
	if err = json.NewDecoder(search.Body).Decode(&result); err != nil {
		return "", err
	}

	hits := result["hits"].(map[string]interface{})
	hitsList := hits["hits"].([]interface{})

	// 获取最新消息ID
	latestID := ""
	if len(hitsList) > 0 {
		latestHit := hitsList[len(hitsList)-1].(map[string]interface{})
		latestID = latestHit["_id"].(string)
	} else {
		return "", errors.New("no messages id found")
	}

	return latestID, nil
}

func GetMessage(index, id string) ([]*Message, error) {
	// 设置查询体
	query := `{"query":{"match":{"session_id":"` + id + `"}},"sort":[{"created_at":{"order":"asc"}}]}`

	// 根据查询体查询出具体的内容
	search, err := _elastic.Search(
		_elastic.Search.WithIndex(index),
		_elastic.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	// 解析消息体
	var result map[string]interface{}
	if err = json.NewDecoder(search.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})
	hitsList := hits["hits"].([]interface{})

	// 遍历hitsList列表
	var MessageList []*Message
	for _, hit := range hitsList {
		hitObject := hit.(map[string]interface{})
		messageId, ok := hitObject["_id"].(string)
		if !ok {
			logger.LogrusObj.Warning("Failed to get message ID")
			continue
		}

		// 访问 _source 字段
		source, ok := hitObject["_source"].(map[string]interface{})
		if !ok {
			logger.LogrusObj.Warning("Failed to get source field or not a map")
			continue
		}

		createdAt, ok := source["created_at"].(float64)
		if !ok {
			logger.LogrusObj.Warning("Failed to get created_at field or not a float64")
			continue
		}

		// 解析 history 字段中的 JSON 字符串
		historyStr, ok := source["history"].(string)
		if !ok {
			logger.LogrusObj.Warning("Failed to get history field or not a string")
			continue
		}
		var history map[string]interface{}
		if err := json.Unmarshal([]byte(historyStr), &history); err != nil {
			logger.LogrusObj.Warning("Failed to decode history JSON:", err)
			continue
		}

		messageType, ok := history["type"].(string)
		if !ok {
			logger.LogrusObj.Warning("Failed to get message type")
			continue
		}
		messageContent, ok := history["data"].(map[string]interface{})["content"].(string)
		if !ok {
			logger.LogrusObj.Warning("Failed to get message content")
			continue
		}

		msg := &Message{
			MessageId: messageId,
			CreatedAt: createdAt,
			Type:      messageType,
			Content:   messageContent,
		}

		MessageList = append(MessageList, msg)
	}

	return MessageList, nil
}
