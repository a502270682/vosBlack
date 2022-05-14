package utils

import (
	"crypto/md5"
	"fmt"
)

func Encrypt(str string) string {
	srcCode := md5.Sum([]byte(str))
	// md5.Sum函数加密后返回的是字节数组，需要转换成16进制形式
	code := fmt.Sprintf("%x", srcCode)
	return code
}
