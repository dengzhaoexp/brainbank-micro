package idmaker

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func InitSnowflakeNode(machineID int64) {
	n, err := snowflake.NewNode(machineID)
	if err != nil {
		panic(err)
	}
	node = n
	return
}

func GenerateNeuronID() string {
	// 根据雪花算法生成id
	snowflakeID := node.Generate().Int64()
	for snowflakeID < 0 {
		snowflakeID = node.Generate().Int64()
	}

	// 构造要哈希的字符串，将Snowflake ID转换为Base62编码
	hashStr := "neuron_" + base62Encode(snowflakeID)

	// 计算哈希值
	hash := hashString(hashStr)

	return hash
}

func GenerateDocumentID() string {
	// 根据雪花算法生成id
	snowflakeID := node.Generate()
	// 返回哈希值
	return snowflakeID.String()
}

func base62Encode(n int64) string {
	if n < 0 {
		return ""
	}

	// 将整数转换为Base62编码字符串
	var base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result string
	for n > 0 {
		result = string(base62Chars[n%62]) + result
		n = n / 62
	}
	return result
}

func hashString(s string) string {
	// 计算字符串的 SHA256 哈希值
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)

	// 对哈希值进行Base64编码
	encoded := base64.StdEncoding.EncodeToString(hashed)

	return encoded
}
