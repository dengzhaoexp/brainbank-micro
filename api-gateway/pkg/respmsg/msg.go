package respmsg

import "api-gateway/pkg/statuscode"

var msgFlag = map[uint32]string{
	statuscode.Success:                  "成功",
	statuscode.Error:                    "失败",
	statuscode.InvalidParams:            "请求参数无效",
	statuscode.EmailAlreadyRegistered:   "邮箱已被注册，请直接登录",
	statuscode.EmailNotRegister:         "用户尚未注册",
	statuscode.UserAccountDisable:       "该账号已被封禁",
	statuscode.NotMatchAccountPwd:       "账户密码不正确",
	statuscode.TokenExpiration:          "用户登录状态失效，请重新登录",
	statuscode.EmailCaptchaNotMatched:   "邮箱验证码错误",
	statuscode.InvalidToken:             "Token无效",
	statuscode.EmailCaptchaExpiration:   "邮箱验证码过期",
	statuscode.NeuronNotFound:           "神经元不存在",
	statuscode.UserHasNoNeuron:          "用户没有该神经元",
	statuscode.NameExisted:              "名称已经存在",
	statuscode.NameSameAsOriginal:       "新名称与旧名称相同",
	statuscode.InvalidDocumentOwnership: "文档所有权或神经元ID不匹配",
	statuscode.DocumentNotFound:         "文件不存在",
	statuscode.FailLoadTemplate:         "加载提示词模版失败",
	statuscode.FailedRequestAI:          "请求AI端失败",
	statuscode.FailedGetMessageID:       "获取消息的ID失败",
	statuscode.NULLConversationID:       "会话ID为空",
	statuscode.AbandonConversation:      "聊天会话已被丢弃",
}

func GetMsg(code uint32) string {
	msg, ok := msgFlag[code]
	if !ok {
		return msgFlag[statuscode.Error]
	}
	return msg
}
