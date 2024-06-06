package handler

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
)

type CaptchaResponse struct {
	CaptchaId string `json:"captchaID"` //验证码Id
	ImageUrl  string `json:"imageUrl"`  //验证码图片url
}

func Captcha() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 创建验证码
		length := captcha.DefaultLen
		captchaId := captcha.NewLen(length)
		var iCaptcha CaptchaResponse
		iCaptcha.CaptchaId = captchaId
		iCaptcha.ImageUrl = "/captchaImg/" + captchaId + ".png"
		ctx.JSON(http.StatusOK, iCaptcha)
	}
}

func CaptchaImg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		w, r := ctx.Writer, ctx.Request
		// 解析url地址
		dir, file := path.Split(r.URL.Path)
		ext := path.Ext(file)
		id := file[:len(file)-len(ext)]
		if ext == "" {
			ext = "png"
		}

		if r.FormValue("reload") != "" {
			captcha.Reload(id)
		}
		download := path.Base(dir) == "download"

		// 设置返回的响应头
		w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// 验证码图片
		var content bytes.Buffer
		switch ext {
		case ".png":
			w.Header().Set("Content-Type", "image/png")
			_ = captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
		case ".wav":
			w.Header().Set("Content-Type", "audio/x-wav")
			_ = captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
		}

		if download {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		// 响应数据
		http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	}
}
