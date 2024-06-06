package core

import (
	"context"
	"fmt"
	"github.com/go-mail/mail"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
	"user/config"
	"user/pkg/consts"
	"user/pkg/statuscode"
	"user/pkg/utils/captcha"
	"user/pkg/utils/idmaker"
	iLogger "user/pkg/utils/logger"
	"user/repositry/cache"
	"user/repositry/dao"
	model2 "user/repositry/model"
	"user/userService"
)

type UserService struct {
}

func (s *UserService) UserRegister(ctx context.Context, req *userService.RegisterRequest, resp *userService.UserServiceResponse) error {
	// 默认状态为成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 创建操作数据库的对象
	userDao := dao.NewUserDao(ctx)

	// 检测用户输入的邮箱是否已经被注册
	exists, err := userDao.IsEmailExists(req.EmailAddress)
	if exists {
		// 用户已经存在
		code = statuscode.EmailAlreadyRegistered
		// 返回体
		resp.Code = uint32(code)
		return nil
	}
	if err != nil {
		// 根据邮箱查询时发生未知错误，但是不影响后续注册
		iLogger.LogrusObj.Error("用户注册时查询用户邮件时发生错误:", err)
	}

	// 校验邮件验证码
	code, err = VerifyEmailCodeFromCache(req.EmailAddress, 0, req.EmailCaptcha)
	if err != nil {
		// 校验过程发生错误
		iLogger.LogrusObj.Error("用户注册时获取邮件验证码发生错误:", err)
		return err
	}
	if code != statuscode.Success {
		// 校验失败
		resp.Code = uint32(code)
		return nil
	}

	// 邮箱校验通过，注册用户到数据库
	userId := idmaker.GenerateUserID()
	user := &model2.User{
		UserId:       userId,
		UserName:     req.EmailAddress,
		Nickname:     req.UserNickname,
		Identity:     "regular user",
		EmailAddress: req.EmailAddress,
		Status:       "active",
		TotalSpace:   cast.ToInt64(consts.InitStorageSpace) * 1024 * 1024,
		UsedSpace:    0,
	}

	// 用户密码加密保存
	err = user.SetPassword(req.AccountPassword)
	if err != nil {
		iLogger.LogrusObj.Error("用户注册时加密用户密码失败:", err)
		return err
	}

	// 用户数据入库
	err = userDao.AddUser(user)
	if err != nil {
		iLogger.LogrusObj.Error("用户注册时新用户入库失败:", err)
		return err
	}

	return nil
}

func (s *UserService) UserLogin(ctx context.Context, req *userService.LoginRequest, resp *userService.UserLoginResponse) error {
	// 状态默认为成功
	code := statuscode.Success
	resp.Code = uint32(code)
	resp.Data = nil

	// 创建操作用户数据库的对象
	userDao := dao.NewUserDao(ctx)

	// 检测是否为注册用户
	exits, err := userDao.IsEmailExists(req.EmailAddress)
	if err != nil {
		iLogger.LogrusObj.Error("用户登录时查询用户邮件时发生错误:", err)
		return err
	}
	if !exits {
		// 用户不存在，提醒用户先注册再登录
		code = statuscode.EmailNotRegister
		resp.Code = uint32(code)
		return nil
	}

	// 获取该注册用户
	user, err := userDao.GetUserByEmail(req.EmailAddress)
	if err != nil {
		iLogger.LogrusObj.Error("用户登录时根据邮件地址获取用户时发生错误:", err)
		return err
	}

	// 检验账户状态
	if strings.EqualFold(user.Status, "baned") {
		code = statuscode.UserAccountDisable
		resp.Code = uint32(code)
		return nil
	}

	// 校验用户密码
	if !user.CheckPassword(req.AccountPassword) {
		// 用户密码错误
		code = statuscode.NotMatchAccountPwd
		resp.Code = uint32(code)
		return nil
	}

	// 返回数据
	data := &userService.UserData{
		UserId:   user.UserId,
		Nickname: user.Nickname,
		Identity: user.Identity,
	}
	resp.Data = data
	return nil
}

func (s *UserService) UserStorage(ctx context.Context, req *userService.StorageRequest, resp *userService.StorageResponse) error {
	// 状态默认为成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 操作数据库对象
	userDao := dao.NewUserDao(ctx)

	// 根据用户id查询用户
	user, err := userDao.GetUserByUserID(req.UserId)
	if err != nil {
		iLogger.LogrusObj.Error("获取容量时根据用户id查询用户时出错:", err)
		return err
	}

	// 返回数据
	data := &userService.StorageData{
		UsedSpace:  user.UsedSpace,
		TotalSpace: user.TotalSpace,
	}

	// 序列化返回数据
	resp.Data = data
	return nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *userService.UpdatePwdRequest, resp *userService.UserServiceResponse) error {
	// 状态默认为成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 创建操作数据库的对象
	userDao := dao.NewUserDao(ctx)

	// 获取用户
	user, err := userDao.GetUserByUserID(req.UserId)
	if err != nil {
		iLogger.LogrusObj.Error("更新用户密码时根据用户id查询用户时出错:", err)
		return err
	}

	// 检查旧密码是否正确
	if !user.CheckPassword(req.AccountPassword) {
		// 旧密码错误
		code = statuscode.NotMatchAccountPwd
		resp.Code = uint32(code)
		return nil
	}

	// 更新用户密码
	err = user.SetPassword(req.NewPassword)
	if err != nil {
		iLogger.LogrusObj.Error("更新用户密码时更新密码时出错:", err)
		return err
	}

	// 用户更新
	err = userDao.UpdateUser(user)
	if err != nil {
		iLogger.LogrusObj.Error("更新用户密码时用户入库时出错:", err)
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *userService.ResetPwdRequest, resp *userService.UserServiceResponse) error {
	// 请求默认成功
	code := statuscode.Success
	resp.Code = uint32(code)

	userDao := dao.NewUserDao(ctx)
	// 检测用户是否存在
	isExists, err := userDao.IsEmailExists(req.EmailAddress)
	if err != nil {
		iLogger.LogrusObj.Error("重置用户密码时查询用户邮件时发生错误:", err)
		return err
	}

	// 用户不存在
	if !isExists {
		code = statuscode.EmailNotRegister
		resp.Code = uint32(code)
		return nil
	}

	// 校验邮件验证码
	code, err = VerifyEmailCodeFromCache(req.EmailAddress, 1, req.EmailCaptcha)
	if err != nil {
		// 校验过程发生错误
		iLogger.LogrusObj.Error("重置用户密码时获取邮件验证码发生错误:", err)
		return err
	}
	if code != statuscode.Success {
		// 校验失败
		resp.Code = uint32(code)
		return nil
	}

	// 校验成功，通过邮件获取user
	user, err := userDao.GetUserByEmail(req.EmailAddress)
	if err != nil {
		iLogger.LogrusObj.Error("重置用户密码时通过邮件获取用户发生错误:", err)
		return err
	}

	// 更新密码
	if err = user.SetPassword(req.AccountPassword); err != nil {
		iLogger.LogrusObj.Error("重置用户密码时设置用户密码出错:", err)
		return err
	}

	// 更新用户入库
	if err = userDao.UpdateUser(user); err != nil {
		iLogger.LogrusObj.Error("重置用户密码时候更新用户出错:", err)
		return err
	}
	return nil
}

func (s *UserService) SendEmail(ctx context.Context, req *userService.SendEmailRequest, resp *userService.UserServiceResponse) error {
	// 默认状态为成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 创建操作数据库的对象
	userDao := dao.NewUserDao(ctx)
	isExists, err := userDao.IsEmailExists(req.EmailAddress)
	if err != nil {
		iLogger.LogrusObj.Error("查询用户邮件是否存在时出错:", err)
		return err
	}

	// 操作邮件的对象
	mailDao := dao.NewMailDaoByDb(userDao.DB)

	// 根据type获取title
	m := new(model2.Mail)
	var key string
	switch req.Type {
	case 0: // 用户注册逻辑
		if isExists {
			// 用户邮箱已被注册
			code = statuscode.EmailAlreadyRegistered
			resp.Code = uint32(code)
			return nil
		}
		// 获取注册邮件内容
		m, err = mailDao.GetResource(consts.MailRegisterTitle)
		if err != nil {
			iLogger.LogrusObj.Error("查询注册邮件title出错:", err)
			return err
		}

		// 缓存的key
		key = cache.VerificationCodeCacheKey(int(req.Type), req.EmailAddress)

	case 1: // 修改密码逻辑
		if !isExists {
			// 邮箱尚未注册
			code = statuscode.EmailNotRegister
			resp.Code = uint32(code)
			return nil
		}

		// 获取修改密码邮件内容
		m, err = mailDao.GetResource(consts.MailResetPwdTitle)
		if err != nil {
			iLogger.LogrusObj.Error("查询修改密码邮件title出错:", err)
			return err
		}

		// 缓存的key
		key = cache.VerificationCodeCacheKey(int(req.Type), req.EmailAddress)
	}

	// 发送邮件的内容
	content, err := mailDao.GetResource(consts.MailContentId)
	if err != nil {
		iLogger.LogrusObj.Error("查询注册邮件content出错:", err)
		return err
	}

	// 生成验证码并生成邮件内容
	c := captcha.GenerateEmailCaptcha()
	mailText := fmt.Sprintf(m.Setting+"\n"+content.Setting, c, strconv.Itoa(consts.EmailCaptchaExpiration))

	// 获取redis客户端并设置缓存
	rClient := cache.GetRedisClient()
	expiration := consts.EmailCaptchaExpiration * time.Minute
	err = rClient.Set(key, c, expiration).Err()
	if err != nil {
		iLogger.LogrusObj.Error("设置邮件验证码缓存时出错:", err)
		return err
	}

	// 发送邮件
	eConfig := config.Config.Email
	msg := mail.NewMessage()
	msg.SetHeader("From", eConfig.SmtpEmail)
	msg.SetHeader("To", req.EmailAddress)
	msg.SetHeader("Subject", "BrainBank")
	msg.SetBody("text/plain", mailText)

	// 创建拨号器
	d := mail.NewDialer(eConfig.SmtpHost, 465, eConfig.SmtpEmail, eConfig.SmtpPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 发送
	if err := d.DialAndSend(msg); err != nil {
		iLogger.LogrusObj.Error("发送邮件时出错:", err)
		return err
	}

	return nil
}
