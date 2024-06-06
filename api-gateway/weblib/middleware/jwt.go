package middleware

import (
	"api-gateway/pkg/respmsg"
	"api-gateway/pkg/statuscode"
	"api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取token
		authHeader := ctx.GetHeader("Authorization")
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		// 判断是否携带token
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":   statuscode.InvalidParams,
				"status": "请求失败",
				"info":   "需要提供有效的 Token",
			})
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil || claims == nil {
			code := statuscode.InvalidToken
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":   code,
				"status": "请求失败",
				"info":   respmsg.GetMsg(uint32(code)),
			})
			ctx.Abort()
			return
		}

		// 判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt {
			code := statuscode.TokenExpiration
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":   code,
				"status": "请求失败",
				"info":   respmsg.GetMsg(uint32(code)),
			})
			ctx.Abort()
			return
		}

		// 携带数据
		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
