package core

import (
	"github.com/go-redis/redis"
	"strings"
	"user/pkg/statuscode"
	"user/repositry/cache"
)

func VerifyEmailCodeFromCache(emailAddress string, emailType int, captcha string) (int, error) {
	// 获取缓存的key
	key := cache.VerificationCodeCacheKey(emailType, emailAddress)

	// 从缓存中获取验证码
	client := cache.GetRedisClient()
	val, err := client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return statuscode.EmailCaptchaExpiration, nil
		}
		return statuscode.Error, err
	}

	// 验证码校验失败
	if !strings.EqualFold(val, captcha) {
		return statuscode.EmailCaptchaNotMatched, nil
	}

	// 验证码校验正确
	return statuscode.Success, nil
}
