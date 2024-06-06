package weblib

import (
	"api-gateway/weblib/handler"
	"api-gateway/weblib/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	// 初始化
	ginRouter := gin.Default()
	// 加载中间件
	ginRouter.Use(middleware.Cors(),
		middleware.Cors(),
		middleware.InitMiddleware(service),
		middleware.ErrorMiddleware())
	v1 := ginRouter.Group("/")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})

		// 验证码
		v1.GET("captcha", handler.Captcha())
		v1.GET("captchaImg/:img", handler.CaptchaImg())
		v1.POST("sendEmail", handler.SendEmail())

		// 用户服务
		v1.POST("user/register", handler.UserRegister())
		v1.POST("user/login", handler.UserLogin())
		v1.POST("user/resetPwd", handler.ResetPwd())

		// 登录保护
		authed := v1.Group("")
		authed.Use(middleware.JWT())
		{
			// 用户
			authed.POST("user/storage", handler.UserStorage())
			authed.POST("user/updatePwd", handler.UpdatePwd())

			// 问答
			authed.POST("conversation", handler.Conversation())
			authed.GET("conversations", handler.Conversations())
			authed.GET("conversation/:conversationId", handler.GetConversation())

			// 文件
			authed.GET("getAvatar/:userId", handler.GetAvatar())
			authed.POST("updateAvatar", handler.UpdateAvatar())
			// Neuron
			authed.PUT("neuron", handler.CreateNeuron())
			authed.DELETE("neuron/:neuronId", handler.DeleteNeuron())
			authed.POST("neuron/rename", handler.RenameNeuron())
			authed.GET("neuron", handler.ListNeuron())
			// Document
			authed.POST("document/upload", handler.UploadDocument())
			authed.POST("document/list", handler.ListDocuments())
			authed.POST("document/rename", handler.RenameDocument())
			authed.DELETE("document/delete", handler.DeleteDocument())
			authed.GET("document/bin", handler.ListDocumentInBin())
			authed.POST("document/recovery", handler.RecoveryDocument())
		}
	}
	return ginRouter
}
