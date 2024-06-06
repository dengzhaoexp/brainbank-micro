package handler

import (
	"api-gateway/pkg/respmsg"
	"api-gateway/pkg/statuscode"
	"api-gateway/pkg/utils"
	"api-gateway/pkg/utils/logger"
	"api-gateway/service/userService"
	"api-gateway/types/userdtos"
	"errors"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 参数绑定与校验
		var req userdtos.UserRegisterRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 校验验证码
		if !captcha.VerifyString(req.CaptchaId, req.Captcha) {
			// 未通过验证码校验
			returnError(ctx, http.StatusOK, errors.New("验证码错误"))
			return
		}

		// 绑定请求参数
		rpcReq := &userService.RegisterRequest{
			EmailAddress:    req.EmailAddress,
			UserNickname:    req.UserNickname,
			AccountPassword: req.AccountPassword,
			EmailCaptcha:    req.EmailCaptcha,
		}

		// 从gin中获取服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		rpcResp, err := service.UserRegister(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 绑定返回数据
		resp := userdtos.UserServiceResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
		}

		// 成功处理
		ctx.JSON(http.StatusOK, resp)
	}
}

func UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 校验绑定与校验
		var req userdtos.UserLoginRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 校验验证码:
		//if !captcha.VerifyString(req.CaptchaId, req.Captcha) {
		//	// 未通过验证码校验
		//	returnError(ctx, http.StatusOK, errors.New("验证码错误"))
		//	return
		//}

		// 绑定参数
		rpcReq := &userService.LoginRequest{
			EmailAddress:    req.EmailAddress,
			AccountPassword: req.AccountPassword,
		}

		// 从gin中取出服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		rpcResp, err := service.UserLogin(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 返回体
		resp := &userdtos.UserLoginResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}

		// 校验rpcResp.Code
		if rpcResp.Code != statuscode.Success {
			resp.Status = "请求失败"
			ctx.JSON(http.StatusOK, resp)
			return
		}

		// 生成token
		token, err := utils.GenerateToken(rpcResp.Data.UserId, rpcResp.Data.Identity)
		if err != nil {
			logger.LogrusObj.Error("生成token失败:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 返回状态
		resp.Token = token
		ctx.JSON(http.StatusOK, resp)
	}
}

func UserStorage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从ctx中获取用户id
		userId := ctx.Keys["userId"].(string)

		// 参数绑定
		rpcReq := &userService.StorageRequest{
			UserId: userId,
		}

		// 从gin中获取服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		rpcResp, err := service.UserStorage(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 绑定返回数据
		resp := &userdtos.StorageResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}

		// 成功处理
		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdatePwd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		var req userdtos.UpdatePwdRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 参数绑定
		rpcReq := &userService.UpdatePwdRequest{
			UserId:          ctx.Keys["userId"].(string),
			AccountPassword: req.AccountPassword,
			NewPassword:     req.NewPassword,
		}

		// 从gin中获取服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		rpcResp, err := service.UpdatePassword(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 绑定返回数据
		resp := userdtos.UserServiceResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
		}

		// 成功处理
		ctx.JSON(http.StatusOK, resp)
	}
}

func ResetPwd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		var req userdtos.ResetPwdRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 参数绑定
		rpcReq := &userService.ResetPwdRequest{
			EmailAddress:    req.EmailAddress,
			AccountPassword: req.AccountPassword,
			EmailCaptcha:    req.EmailCaptcha,
		}

		// 从gin中获取服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		rpcResp, err := service.ResetPassword(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 绑定返回数据
		resp := userdtos.UserServiceResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
		}

		// 成功处理
		ctx.JSON(http.StatusOK, resp)
	}
}

func SendEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 参数绑定与校验
		var req userdtos.SendEmailRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 校验验证码
		if !captcha.VerifyString(req.CaptchaId, req.Captcha) {
			// 未通过验证码校验
			returnError(ctx, http.StatusOK, errors.New("验证码错误"))
			return
		}

		// 绑定请求参数
		t, err := strconv.Atoi(req.Type)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}
		rpcReq := &userService.SendEmailRequest{
			EmailAddress: req.EmailAddress,
			Type:         uint32(t),
		}

		// 从gin中获取服务实例
		service := ctx.Keys["rpcUserService"].(userService.UserService)
		resp, err := service.SendEmail(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 成功处理
		ctx.JSON(http.StatusOK, resp)
	}
}
