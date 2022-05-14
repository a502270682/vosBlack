package utils

import (
	"fmt"
	"testing"
)

func TestEncryot(t *testing.T) {
	str := "123456789"
	sign := Encrypt(str)
	fmt.Println(sign)
}
