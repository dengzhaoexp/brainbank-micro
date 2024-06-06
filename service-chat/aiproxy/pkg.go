package aiproxy

import (
	"bytes"
	"net/http"
)

func postRequest(url string, requestBody []byte) (*http.Response, error) {
	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发起请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
