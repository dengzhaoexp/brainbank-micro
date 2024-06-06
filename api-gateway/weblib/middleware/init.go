package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// InitMiddleware 接受服务实例，并存在gin.Key中
func InitMiddleware(service []interface{}) gin.HandlerFunc {
	// 将服务实例存在gin.keys中
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]interface{})
		ctx.Keys["rpcUserService"] = service[0]
		ctx.Keys["rpcFileService"] = service[1]
		ctx.Keys["rpcChatService"] = service[2]
		ctx.Next()
	}
}

// ErrorMiddleware 错误处理的中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})

			}
		}()
		ctx.Next()
	}
}
