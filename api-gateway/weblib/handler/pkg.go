package handler

import (
	iLogger "api-gateway/pkg/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bindAndValidate 绑定并验证请求参数
func bindAndValidate(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBind(req); err != nil {
		return fmt.Errorf("请求参数错误: %w", err)
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		// 验证未通过,打印错误信息
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			iLogger.LogrusObj.Error("Field: %s, Tag: %s, Error: %s\n", err.Field(), err.Tag(), err.ActualTag())
		}
		return fmt.Errorf("缺少必要参数: %w", err)
	}

	return nil
}

// returnError 返回错误信息
func returnError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, gin.H{"error": err.Error(), "status": "请求失败"})
}
