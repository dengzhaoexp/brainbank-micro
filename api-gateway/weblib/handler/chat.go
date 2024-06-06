package handler

import (
	"api-gateway/pkg/respmsg"
	"api-gateway/pkg/statuscode"
	iLogger "api-gateway/pkg/utils/logger"
	"api-gateway/service/chatService"
	"api-gateway/types/chatdtos"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func Conversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 参数绑定与校验
		var req chatdtos.ConversationRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 获取rpc服务
		service, ok := ctx.Keys["rpcChatService"].(chatService.ConversationService)
		if !ok {
			returnError(ctx, http.StatusInternalServerError, errors.New("rpcChatService.Client not found in context"))
			return
		}

		// 创建标题的请求体
		titleReq := chatService.CreateTitleRequest{
			UserId:         ctx.Keys["userId"].(string),
			MessageContent: req.Message.Content.Parts,
			Language:       req.Language,
		}

		// 判断是否为聊天的第一次请求
		var titleResp *chatService.CreateTitleResponse
		if req.ConversionID == "" {
			// 调用流式服务
			titleStream, err := service.CreateTitle(ctx, &titleReq)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}

			// 循环监听消息
			for {
				// 阻塞接受响应
				titleResp, err = titleStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					// 非正常结束，报500错误
					ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
					return
				}

				req.ConversionID = titleResp.ConversionId
				iLogger.LogrusObj.Info(titleResp.TitleContent)

				// 判断请求失败
				if titleResp.Code != statuscode.Success {
					ctx.JSON(http.StatusOK, gin.H{"code": titleResp.Code, "info": respmsg.GetMsg(titleResp.Code)})
					return
				}

				break
			}
		}

		// 消息请求体
		msgReq := chatService.AssistantMessageRequest{
			ConversionId:     req.ConversionID,
			ConversationMode: req.ConversationMode.Kind,
			Content:          req.Message.Content.Parts,
			ContentType:      "text",
			Language:         req.Language,
		}

		// 获取消息流式服务
		messageStream, err := service.AssistantMessage(ctx, &msgReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 获取消息响应
		messageResp, err := messageStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			// 非正常结束，报500错误
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		if messageResp.Code != statuscode.Success {
			ctx.JSON(http.StatusOK, gin.H{"code": messageResp.Code, "info": respmsg.GetMsg(messageResp.Code)})
			return
		}
		iLogger.LogrusObj.Info(messageResp.Messages)

		// 获取消息成功，以SSE发送
		SendMessageToClient(ctx, &req, titleResp, messageResp)
	}
}

func Conversations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		offsetStr := ctx.Query("offset")
		limitStr := ctx.Query("limit")
		order := ctx.Query("order") // updated

		// 校验参数
		if offsetStr == "" || limitStr == "" || order == "" {
			returnError(ctx, http.StatusBadRequest, errors.New("can't ignore offset or limit or order"))
			return
		}

		// 转换格式
		offset, err := strconv.Atoi(offsetStr)
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			returnError(ctx, http.StatusBadRequest, errors.New("error format in offset or limit"))
			return
		}

		// 创建rpc服务请求体
		rpcReq := chatService.ConversationsRequest{
			UserId: ctx.Keys["userId"].(string),
			Offset: uint32(offset),
			Limit:  uint32(limit),
			Order:  order,
		}

		// 获取服务
		rpcService, ok := ctx.Keys["rpcChatService"].(chatService.ConversationService)
		if !ok {
			returnError(ctx, http.StatusInternalServerError, errors.New("rpcChatService.Client not found in context"))
			return
		}

		// 调用具体服务
		rpcResp, err := rpcService.Conversations(ctx, &rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 判断业务是否成功
		if rpcResp.Code != statuscode.Success {
			ctx.JSON(http.StatusOK, gin.H{"code": rpcResp.Code, "msg": respmsg.GetMsg(rpcResp.Code)})
		}

		// 绑定响应数据
		cItems := make([]*chatdtos.ConversationItem, len(rpcResp.Items))
		for i, _ := range rpcResp.Items {
			c := rpcResp.Items[i]
			item := &chatdtos.ConversationItem{
				ConversationID:         c.ConversationId,
				Title:                  c.ConversationTitle,
				CreateTime:             c.CreateTime,
				UpdateTime:             c.UpdateTime,
				Mapping:                c.Mapping,
				CurrentNode:            c.CurrentNode,
				ConversationTemplateId: c.ConversationTemplateId,
				GizmoId:                c.GizmoId,
				IsArchived:             c.IsArchived,
				WorkspaceId:            c.WorkspaceId,
			}
			cItems[i] = item
		}

		// 绑定响应数据
		resp := &chatdtos.ConversationsResponse{
			Items:                   cItems,
			Total:                   int(rpcResp.Total),
			Limit:                   int(rpcResp.Limit),
			Offset:                  offset,
			HasMissingConversations: rpcResp.HasMissingConversations,
		}
		ctx.JSON(http.StatusOK, resp)

	}
}

func GetConversation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取conversationId参数的值
		conversationId := ctx.Param("conversationId")

		// 获取rpc服务
		service, ok := ctx.Keys["rpcChatService"].(chatService.ConversationService)
		if !ok {
			returnError(ctx, http.StatusInternalServerError, errors.New("rpcChatService.Client not found in context"))
			return
		}

		// 获取聊天消息的请求体
		rpcReq := chatService.GetConversationRequest{ConversationId: conversationId}

		// 调用具体的服务函数
		rpcResp, err := service.GetConversation(ctx, &rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 解析响应数据
		var mapping chatdtos.Mapping
		for k, v := range rpcResp.Mapping {
			msg := chatdtos.MessageStruct{
				ID:         v.MessageId,
				CreateTime: v.CreateTime,
				Content: chatdtos.Content{
					ContentType: "text",
					Parts:       []string{v.Content},
				},
				Type: v.Type,
			}

			mapping.Mapping[k] = msg
		}

		// 响应数据
		data := chatdtos.GetConversationResponse{
			ConversationID: conversationId,
			Title:          rpcResp.ConversationTitle,
			CreateTime:     rpcResp.CreatedTime,
			UpdateTime:     rpcResp.UpdateTime,
			Mapping:        mapping,
		}

		ctx.JSON(http.StatusOK, data)
	}

}

// SendMessageToClient 将消息发送给客户端
func SendMessageToClient(ctx *gin.Context, req *chatdtos.ConversationRequest, title *chatService.CreateTitleResponse, msg *chatService.AssistantMessageResponse) {
	// 设置响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	parts := ""
	for _, msg := range msg.Messages {
		if req.Language == "zh-CN" {
			// 默认按照英语处理方式
			parts = parts + msg
		} else {
			// 默认按照英语处理方式
			parts = parts + " " + msg
		}

		msgObj := chatdtos.ConversationResponse{
			Message: chatdtos.Message{
				Author: chatdtos.Author{Role: req.Message.Author.Role},
				Content: chatdtos.Content{
					ContentType: req.Message.Content.ContentType,
					Parts:       []string{parts},
				},
			},
			ConversationID: req.ConversionID,
		}

		// 将对象转换为 JSON
		marshal, err := json.Marshal(msgObj)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 将 JSON 格式的数据发送给客户端
		ctx.SSEvent("message", string(marshal))
	}

	if title != nil {
		titleObj := chatdtos.TitleResponse{
			Type:           "title_generation",
			Title:          title.TitleContent,
			ConversationID: title.ConversionId,
		}

		mTitle, err := json.Marshal(titleObj)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		// 将 JSON 格式的数据发送给客户端
		ctx.SSEvent("message", string(mTitle))
	}

	ctx.SSEvent("message", "[Done]")
}
