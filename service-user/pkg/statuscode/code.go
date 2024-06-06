package statuscode

const (
	Success = 200
	Error   = 500

	EmailAlreadyRegistered = 10001
	EmailNotRegister       = 10002
	UserAccountDisable     = 10003
	NotMatchAccountPwd     = 10004
	EmailCaptchaNotMatched = 10006
	EmailCaptchaExpiration = 10007
)
