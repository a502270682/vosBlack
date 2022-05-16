package service

import (
	"context"
	"testing"
)

func TestCommonCheck(t *testing.T) {
	ctx := context.Background()
	realCaller := ""
	enID := 1
	ipID := 1
	callID, caller, callee := "", "", ""
	_ = CommonCheck(ctx, realCaller, enID, ipID, callID, caller, callee)
}
