package statuscode

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	EmailAlreadyRegistered   = 10001
	EmailNotRegister         = 10002
	UserAccountDisable       = 10003
	NotMatchAccountPwd       = 10004
	InvalidToken             = 10005
	EmailCaptchaNotMatched   = 10006
	EmailCaptchaExpiration   = 10007
	NeuronNotFound           = 10008
	UserHasNoNeuron          = 10009
	NameExisted              = 10010
	NameSameAsOriginal       = 10011
	InvalidDocumentOwnership = 10012
	DocumentNotFound         = 10013

	FailLoadTemplate    = 10014
	FailedRequestAI     = 10015
	FailedGetMessageID  = 10016
	NULLConversationID  = 10017
	AbandonConversation = 10018

	TokenExpiration = 20001
)
