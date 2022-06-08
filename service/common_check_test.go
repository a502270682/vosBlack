package service

import (
	"context"
	"testing"
	"vosBlack/common"
)

func TestCommonCheckForMobileBlack(t *testing.T) {
	ctx := context.Background()
	prefix := "806"
	realCaller := "15190163799"
	enID := 1
	ipID := 1
	callID, caller, callee := "", "", ""
	// mobile_black check
	//retStatus := CommonCheck(ctx, realCaller, enID, ipID, callID, caller, callee)
	//if retStatus != common.BlackMobile {
	//	t.Fatal("should get 12001, but get", retStatus)
	//}
	// fq_control check
	phoneType := 1
	retStatus := CommonCheck(ctx, prefix, realCaller, enID, ipID, callID, caller, callee, phoneType)
	if retStatus != common.OutOfFrequency {
		t.Fatal("should get 12003, but get", retStatus)
	}
	t.Log("success")
}
