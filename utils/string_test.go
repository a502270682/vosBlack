package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"vosBlack/proto"
)

func Test_GetPhone(t *testing.T) {
	callee := "8021213810507903"
	prefix, realCallee, phoneType := GetPhone(callee)
	fmt.Printf("prefix: %s realCallee : %s phoneType : %d \n", prefix, realCallee, phoneType)
}

func Test_Unmarshal(t *testing.T) {
	str := "{\"callId\":\"13810507030111\",\"code\":1,\"msg\":\"success\",\"list\":[{\"mobile\":\"13810507133\",\"forbid\":2,\"msg\":\"\\u8d85\\u9891\\u963b\\u6b62\"}],\"transactionId\":\"5214034334\"}"
	rsp := &proto.BlackDongYunRsp{}
	err := json.Unmarshal([]byte(str), rsp)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
	}
	fmt.Printf("%+v", rsp)
}
