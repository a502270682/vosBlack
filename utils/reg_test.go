package utils

import (
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"testing"
)

func Test_FindStringSubmatch(t *testing.T) {
	str := "71615121603408"
	matches := FindStringSubmatch(str)
	fmt.Printf("%+v", matches)
}

func Test_Regxp(t *testing.T) {
	a := pcre.MustCompile("^\\d+(\\d)\\1{2}(?!\\1)\\d$", 0)
	realCallee := "15121601113"
	if len(a.FindIndex([]byte(realCallee), 0)) > 0 {
		fmt.Println("匹配成功")
	} else {
		fmt.Println("匹配失败")
	}
}
