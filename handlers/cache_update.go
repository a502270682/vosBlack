package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"vosBlack/adapter/error_code"
	"vosBlack/adapter/log"
	"vosBlack/adapter/logic"
	"vosBlack/model"
	"vosBlack/proto"
)

func CacheUpdate(ctx context.Context, req *proto.CacheUpdateReq, rsp *proto.CacheUpdateRsp) *error_code.ReplyError {
	// update t_enterprise_blacklist
	enterpriseBlackLists, err := model.GetEnterpriseBlacklistImpl().GetAllEnterpriseBlacklist(ctx)
	if err != nil {
		log.Warnf(ctx, "GetAllEnterpriseBlacklist failed, err:%s", err.Error())
	}
	for _, e := range enterpriseBlackLists {
		err = logic.SetEnterpriseBlackListCache(ctx, e.EnIPID, e.Qianzhui, e)
		if err != nil {
			log.Warnf(ctx, "SetEnterpriseBlackListCache failed, err:%s", err.Error())
			continue
		}
	}
	// update t_enterprise_calllist
	enterpriseCallLists, err := model.GetEnterpriseCalltimelistImpl().GetAll(ctx)
	if err != nil {
		log.Warnf(ctx, "GetEnterpriseCalltimelistImpl failed, err:%s", err.Error())
	}
	callMap := map[string][]*model.EnterpriseCalltimelist{}
	for _, e := range enterpriseCallLists {
		key := fmt.Sprintf("%d-%d", e.EnID, e.BlackID)
		callMap[key] = append(callMap[key], e)
	}
	for key, val := range callMap {
		arr := strings.Split(key, "-")
		enID, _ := strconv.Atoi(arr[0])
		blackID, _ := strconv.Atoi(arr[1])
		err = logic.SetEnterpriseCallTimeListCache(ctx, enID, blackID, val)
		if err != nil {
			log.Warnf(ctx, "SetEnterpriseCallTimeListCache failed, err:%s", err.Error())
			continue
		}
	}
	// update mobile_whitenum
	mws, err := model.GetMobileWhitenumImpl().GetAllActiveMw(ctx)
	if err != nil {
		log.Warnf(ctx, "GetAllActiveMw failed, err:%s", err.Error())
	}
	for _, mw := range mws {
		err = logic.SetMobileWhiteNumCache(ctx, mw.EnID, mw.WhiteNum, mw)
		if err != nil {
			log.Warnf(ctx, "SetMobileWhiteNumCache failed, err:%s", err.Error())
			continue
		}
	}
	// update mobile_pattern
	mps, err := model.GetMobilePatternImpl().GetAllActiveMbLevels(ctx)
	if err != nil {
		log.Warnf(ctx, "GetAllActiveMbLevels failed, err:%s", err.Error())
	}
	mbLevel2MobilePattern := make(map[int][]*model.MobilePattern)
	for _, mp := range mps {
		if this, ok := mbLevel2MobilePattern[mp.MbLevel]; ok {
			this = append(this, mp)
			mbLevel2MobilePattern[mp.MbLevel] = this
		} else {
			var newMps []*model.MobilePattern
			newMps = append(newMps, mp)
			mbLevel2MobilePattern[mp.MbLevel] = newMps
		}
	}
	for mbLevel, thisMps := range mbLevel2MobilePattern {
		err = logic.SetMobilePatternCache(ctx, mbLevel, thisMps)
		if err != nil {
			log.Warnf(ctx, "SetMobilePatternCache failed, err:%s", err.Error())
			continue
		}
	}
	return nil
}
