package cache

import "fmt"

func VerificationCodeCacheKey(kind int, email string) string {
	return fmt.Sprintf("VerificationCodeCacheKey:%s:%d", email, kind)
}
