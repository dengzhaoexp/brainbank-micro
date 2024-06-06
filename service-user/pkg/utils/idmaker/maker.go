package idmaker

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
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

func GenerateUserID() string {
	id, _ := uuid.NewUUID()
	return "user-" + id.String()
}
