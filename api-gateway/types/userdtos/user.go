package userdtos

import "api-gateway/service/userService"

type UserServiceResponse struct {
	Code   uint32 `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	Info   string `form:"info" json:"info"`
}

type SendEmailRequest struct {
	EmailAddress string `form:"email_address" json:"email_address" binding:"required"`
	Captcha      string `form:"captcha" json:"captcha" binding:"required"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" binding:"required"`
	Type         string `form:"types" json:"types" binding:"required"`
}

type UserRegisterRequest struct {
	EmailAddress    string `form:"email_address" json:"email_address" binding:"required"`
	UserNickname    string `form:"user_nickname" json:"user_nickname"`
	AccountPassword string `form:"account_password" json:"account_password" binding:"required"`
	EmailCaptcha    string `form:"email_captcha" json:"email_captcha" binding:"required"`
	Captcha         string `form:"captcha" json:"captcha" binding:"required"`
	CaptchaId       string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type UserLoginRequest struct {
	EmailAddress    string `form:"email_address" json:"email_address" binding:"required"`
	AccountPassword string `form:"account_password" json:"account_password" binding:"required"`
	Captcha         string `form:"captcha" json:"captcha" binding:"required"`
	CaptchaId       string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type UserLoginResponse struct {
	Code   uint32                `form:"code" json:"code"`
	Status string                `form:"status" json:"status"`
	Info   string                `form:"info" json:"info"`
	Data   *userService.UserData `form:"data" json:"data"`
	Token  string                `form:"token" json:"token"`
}

type StorageResponse struct {
	Code   uint32                   `form:"code" json:"code"`
	Status string                   `form:"status" json:"status"`
	Info   string                   `form:"info" json:"info"`
	Data   *userService.StorageData `form:"data" json:"data"`
}

type UpdatePwdRequest struct {
	AccountPassword string `form:"account_password" json:"account_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
}

type ResetPwdRequest struct {
	EmailAddress    string `form:"email_address" json:"email_address" binding:"required"`
	AccountPassword string `form:"account_password" json:"account_password" binding:"required"`
	EmailCaptcha    string `form:"email_captcha" json:"email_captcha" binding:"required"`
}
