package utils

import (
	"fmt"
	"testing"
)

func TestEncryot(t *testing.T) {
	str := "1000021234567891123456"
	sign := Encrypt(str)
	fmt.Println(sign)
}
