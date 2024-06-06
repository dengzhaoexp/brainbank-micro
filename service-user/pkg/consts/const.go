package consts

const (
	JWTSecret = "BrainBank_JWTSecret"

	InitStorageSpace = 5 // 初始存储容量

	UserIdLength    = 15
	TokenExpiration = 24

	EmailCaptchaLength     = 5  // 邮箱验证码的长度
	EmailCaptchaExpiration = 40 // 邮箱验证码有效时间（Min）

	MailContentId     = "content"        // 查询邮件内容的id
	MailRegisterTitle = "register_title" // 注册邮件的title
	MailResetPwdTitle = "resetPwd_title" // 重置密码的title
)
