package logic

import (
	"context"
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	// mb_cache
	mbLevel := 0
	enID := 1
	whiteNum := "13990163799"
	ipID := 1
	prefex := "151"
	blackID := 1
	list, err := GetMobilePatternWithCache(ctx, mbLevel)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range list {
		t.Log(a)
	}

	mwn, err := GetMobileWhiteNumWithCache(ctx, enID, whiteNum)
	if err != nil {
		t.Fatal(err)
	}
	if mwn.WhiteNum != whiteNum {
		t.Fatal("whiteNum not right")
	}

	ebl, err := GetEnterpriseBlackListWithCache(ctx, ipID, prefex)
	if err != nil {
		t.Fatal(err)
	}
	if ebl.PatternLevel != 400 {
		t.Fatal(fmt.Sprintf("ebl: should get(%d), but get(%d)", 400, ebl.PatternLevel))
	}

	ecl, err := GetEnterpriseCallTimeListWithCache(ctx, enID, blackID)
	if err != nil {
		t.Fatal(err)
	}
	if ecl.TimeName != "全天" {
		t.Fatal(fmt.Sprintf("ecl: should get(%s), but get(%s)", "全天", ecl.TimeName))
	}
	t.Log("success")
}
