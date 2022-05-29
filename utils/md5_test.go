package utils

import (
	"fmt"
	"testing"
)

func TestEncryot(t *testing.T) {
	str := "admin@DaLian011423530JkQZfYxqB1"
	sign := Encrypt(str)
	fmt.Println(sign)
}
