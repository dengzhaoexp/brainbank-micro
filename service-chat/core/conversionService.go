package core

import (
	"chat/aiproxy"
	"chat/chatService"
	iLogger "chat/pkg/logger"
	"chat/pkg/statuscode"
	"chat/repositry/dao"
	"chat/repositry/es"
	"chat/repositry/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ConversionServer struct {
	// 此处可以注入需要的依赖,如会话存储等
}

func (server *ConversionServer) CreateTitle(ctx context.Context, req *chatService.CreateTitleRequest, stream chatService.ConversationService_CreateTitleStream) (err error) {
	// 默认状态为成功
	resp := chatService.CreateTitleResponse{
		Code:         statuscode.Success,
		TitleContent: "",
		ConversionId: "",
	}

	// 获取用户会话的所有标题：TODO：优化，只获取当天的标题，避免占用过多的token
	existTitleStr := ""
	cDao := dao.NewResourceDao(ctx)
	cs, err := cDao.GetConversationsByUserId(req.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			iLogger.LogrusObj.Error("根据用户id获取所有聊天时发生错误:", err)
			return err
		}
	} else {
		if len(cs) > 0 {
			titles := make([]string, len(cs))
			for i, c := range cs {
				titles[i] = c.Title
			}
			existTitleStr = strings.Join(titles, "|")
		}
	}

	// 构建请求的消息体
	titleRequest := aiproxy.ConversationTitleRequest{
		Content:        strings.Join(req.MessageContent, " "),
		ExcludedTitles: existTitleStr,
	}
	iLogger.LogrusObj.Info("请求创建会话title的请求体已经创建好:", titleRequest)

	// 获取AI响应
	aiResp, err := aiproxy.GetConversationTitleFromLLM(&titleRequest)
	if err != nil {
		iLogger.LogrusObj.Error("从ai端加载消息出错:", err)
		resp.Code = statuscode.FailedRequestAI
		if err = stream.Send(&resp); err != nil {
			return err
		}
	}

	// 生成会话ID
	conversationId := uuid.New().String()

	// 会话信息入库
	conversation := model.Conversation{
		ConversationId: conversationId,
		Title:          aiResp.Content,
		Status:         "created",
		CreatedBy:      req.UserId,
		MessageCount:   0,
		Language:       req.Language,
		TokensConsumed: 0,
	}

	if err = cDao.CreateConversation(conversation); err != nil {
		iLogger.LogrusObj.Error("会话信息入库失败:", err)
		return err
	}

	// 响应消息
	resp.TitleContent = aiResp.Content
	resp.ConversionId = conversationId
	if err = stream.Send(&resp); err != nil {
		return err
	}

	// 发送 io.EOF 表示流结束
	if err = stream.Close(); err != nil {
		return err
	}

	return err
}

func (server *ConversionServer) AssistantMessage(ctx context.Context, req *chatService.AssistantMessageRequest, stream chatService.ConversationService_AssistantMessageStream) (err error) {
	// 默认状态为成功
	resp := chatService.AssistantMessageResponse{
		Code:      statuscode.Success,
		Messages:  nil,
		MessageId: "",
	}

	// 检查会话id
	if req.ConversionId == "" {
		resp.Code = statuscode.NULLConversationID
		if err := stream.Send(&resp); err != nil {
			return err
		}
		if err = stream.Close(); err != nil {
			return err
		}
		return nil
	}

	// 构建请求消息体
	messageRequest := aiproxy.ConversationMessageRequest{
		Content:          strings.Join(req.Content, "|"),
		ConversationMode: req.ConversationMode,
		ConversationID:   req.ConversionId,
	}

	// 发送请求获取对应的AI消息 TODO:这里需要设置超时时间，然后针对超时必须要有正确的处理
	messageResp, err := aiproxy.GetMessageFromLLM(&messageRequest)
	if err != nil {
		iLogger.LogrusObj.Error("从ai端加载消息出错:", err)
		resp.Code = statuscode.FailedRequestAI
		if err = stream.Send(&resp); err != nil {
			return err
		}
	}
	iLogger.LogrusObj.Info(messageResp.Content)

	// 获取messageID
	messageID, err := es.GetMessageID("chat-history", req.ConversionId)
	if err != nil {
		iLogger.LogrusObj.Error("从es中获取消息的ID失败:", err)
		resp.Code = statuscode.FailedGetMessageID
		if err = stream.Send(&resp); err != nil {
			return err
		}
	}
	iLogger.LogrusObj.Info("获取到消息ID:", messageID)

	// 绑定数据并响应
	resp.Messages = messageResp.Content
	resp.MessageId = messageID
	if err = stream.Send(&resp); err != nil {
		return err
	}
	// 发送 io.EOF 表示流结束
	if err = stream.Close(); err != nil {
		return err
	}

	// 构建数据库对象
	mDao := dao.NewResourceDao(ctx)

	// 解析请求体
	c, err := mDao.GetConversationByID(req.ConversionId)
	if err != nil {
		iLogger.LogrusObj.Error("根据id获取对应conversation时发生错误:", err)
		return err
	}

	// 更新会话信息
	newC := model.Conversation{
		ConversationId: c.ConversationId,
		Title:          c.Title,
		Status:         "on-processing",
		CreatedBy:      c.CreatedBy,
		MessageCount:   c.MessageCount + 1,
		Language:       req.Language,
		TokensConsumed: c.TokensConsumed + messageResp.TokenUsage,
	}
	newC.ID = c.ID

	// 新对象入库
	if err = mDao.UpdateConversation(newC); err != nil {
		iLogger.LogrusObj.Error("在更新会话对象时发生错误:", err)
		return err
	}

	return err
}

func (server *ConversionServer) Conversations(ctx context.Context, req *chatService.ConversationsRequest, resp *chatService.ConversationsResponse) error {
	// 默认状态为请求成功
	resp.Code = statuscode.Success

	// 获取操控数据库的对象
	cDao := dao.NewResourceDao(ctx)

	// 获取所有的Conversations
	cs, err := cDao.GetConversationsByUserIdWithOffset(req.UserId, req.Offset, req.Limit)
	if err != nil {
		resp.Code = statuscode.Error
		iLogger.LogrusObj.Error("通过用户的id查询所有的会话出错:", err)
		return err
	}

	// 如果为空
	if len(cs) == 0 {
		resp.Items = nil
		resp.Total = 0
		resp.Limit = req.Limit
		resp.HasMissingConversations = false
		return nil
	}

	// 遍历赋值
	items := make([]*chatService.ConversationsItem, len(cs))
	for i, _ := range cs {
		c := cs[i]
		item := &chatService.ConversationsItem{
			ConversationId:         c.ConversationId,
			ConversationTitle:      c.Title,
			CreateTime:             c.CreatedAt.String(),
			UpdateTime:             c.UpdatedAt.String(),
			Mapping:                "",
			CurrentNode:            "",
			ConversationTemplateId: "",
			GizmoId:                "",
			IsArchived:             false,
			WorkspaceId:            "",
		}
		items[i] = item
	}

	// 绑定返回数据
	resp.HasMissingConversations = false
	resp.Total = uint32(len(cs))
	resp.Limit = req.Limit
	resp.Items = items
	return err
}

func (server *ConversionServer) GetConversation(ctx context.Context, req *chatService.GetConversationRequest, resp *chatService.GetConversationResponse) error {
	// 默认状态为请求成功
	resp.Code = statuscode.Success

	// 获取操作数据库的对象
	cDao := dao.NewResourceDao(ctx)

	// 获取对应的聊天会话对象
	c, err := cDao.GetConversationByID(req.ConversationId)
	if err != nil {
		iLogger.LogrusObj.Error("通过聊天会话id查询会话时发生错误:", err)
		return err
	}

	// 检查该会话的状态
	if c.Status == "Abandoned" {
		resp.Code = statuscode.AbandonConversation
		return nil
	}

	// 获取完整消息记录
	msgs, err := es.GetMessage("chat-history", req.ConversationId)
	if err != nil {
		iLogger.LogrusObj.Error("获取完整聊天消息记录失败:", err)
		return err
	}

	// 处理返回数据
	for _, msg := range msgs {
		createdAtTime := time.Unix(int64(msg.CreatedAt), 0)
		createdAtStr := createdAtTime.Format("2006-01-02 15:04:05")
		// 创建一个Mapping对象并设置其字段值
		mapping := &chatService.Mapping{
			MessageId:  msg.MessageId,
			CreateTime: createdAtStr,
			Content:    msg.Content,
			Type:       msg.Type,
		}

		// 将mapping对象添加到resp的mapping字段中
		resp.Mapping[msg.MessageId] = mapping
	}

	// 绑定返回数据
	resp.ConversationTitle = c.Title
	resp.CreatedTime = c.CreatedAt.String()
	resp.UpdateTime = c.UpdatedAt.String()

	return err
}
