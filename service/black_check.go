package service

import (
	"context"
	"fmt"
	"vosBlack/adapter/http"
	"vosBlack/adapter/log"
	"vosBlack/adapter/logic"
	"vosBlack/common"
	"vosBlack/model"
	"vosBlack/utils"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"gorm.io/gorm"
)

func GetCompanyByIP(ctx context.Context, ip string) (*model.EnterpriseIplist, error) {
	return model.GetEnterpriseIplistImpl().GetOneByIP(ctx, ip)
}

func GetUserByUserID(ctx context.Context, userID string) (*model.EnterpriseUserlist, error) {
	return model.GetEnterpriseUserlistImpl().GetByUserID(ctx, userID)
}

func GetEnterpriseFeelList(ctx context.Context, enID int) (*model.EnterpriseFeeList, error) {
	return model.GetEnterpriseFeeListImpl().GetOneByEnID(ctx, enID)
}

func IsWhiteNum(ctx context.Context, realCallee string, enID int) (bool, error) {
	whiteNum, err := logic.GetMobileWhiteNumWithCache(ctx, enID, realCallee)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// 在白名单内返回正常号码
	if whiteNum != nil {
		return true, nil
	}
	return false, nil
}

func CommonCheck(ctx context.Context, prefix, realCallee string, enID, ipID int, callID interface{}, caller, callee string, phoneType int) common.RespStatus {
	// 根据前缀和ip获取黑名单规则
	blackRule, err := logic.GetEnterpriseBlackListWithCache(ctx, ipID, prefix)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.NotFound
		}
		return common.SystemInternalError
	}
	var mbRequestCount, mbHitCount, wnHitCount, mpRequestCount, mpHitCount, gwRequestCount, gwHitCount, fqRequestCount, fqHitCount int64
	defer func() {
		mbRequestCount = 1
		err = logic.UpsertEnterpriseApplyHourList(ctx, enID, "", mbRequestCount, mbHitCount, wnHitCount, mpRequestCount, mpHitCount, gwRequestCount, gwHitCount, fqRequestCount, fqHitCount)
		if err != nil {
			log.Warnf(ctx, "UpsertEnterpriseApplyHourList failed, err:%+v", err)
		}
		// todo @feiyangguo 目前是先计算是否超过频次，再累计该次，看后续是否需要调整
		if fqRequestCount > 0 {
			err = logic.AddEnterpriseFqCache(ctx, enID, realCallee, utils.GetLastNDay0TimeStamp(0), 1)
			if err != nil {
				log.Warnf(ctx, fmt.Sprintf("AddEnterpriseFqCache fail, err:%s", err.Error()))
			}
		}
	}()
	// 判断白名单
	if blackRule.IsWhitenum == 1 {
		isExist, err := IsWhiteNum(ctx, realCallee, enID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if isExist {
			wnHitCount = 1
			//logic.UpsertEnterpriseApplyHourList(ctx, enID, "", 1, 0, 1, 0, 0, 0, 0, 0, 0)
			return common.StatusOK
		}
	}
	// 判断呼叫时间段
	if blackRule.IsCalltime == 1 {
		callTimeList, err := logic.GetEnterpriseCallTimeListWithCache(ctx, enID, blackRule.NID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if err == nil && callTimeList != nil {
			if !utils.IsReachTime(callTimeList.BeginHour, callTimeList.BeginMinute, callTimeList.EndHour, callTimeList.Edminute) {
				return common.UnReachTime
			}
		}
	}
	// 判断靓号
	if blackRule.PatternLevel != -1 {
		mobilePatternList, err := logic.GetMobilePatternWithCache(ctx, blackRule.PatternLevel)
		if err != nil {
			return common.SystemInternalError
		}
		mpRequestCount = 1
		for _, value := range mobilePatternList {
			reg := pcre.MustCompile(value.Pattern, 0)
			if len(reg.FindIndex([]byte(realCallee), 0)) > 0 {
				mpHitCount = 1
				return common.PrettyNumber
			}
		}
	}
	// 判断本地黑名单
	if blackRule.BlacknumLevel != -1 {
		tablePrefix := ""
		if phoneType == 0 {
			tablePrefix = "0"
		} else {
			tablePrefix = realCallee[:3]
		}
		mobile, err := model.GetMobileBlackApi().GetOneByMobileAll(ctx, tablePrefix, realCallee)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if mobile != nil && mobile.MobileAll == realCallee {
			mbHitCount = 1
			return common.BlackMobile
		}
	}
	if blackRule.IsFrequency == 1 {
		if blackRule.CallCycle != -1 && blackRule.CallCount > 0 {
			startDate := utils.GetLastNDay0TimeStamp(blackRule.CallCycle)
			frequencyCount, err := logic.GetEnterpriseFqFromStartDay(ctx, enID, realCallee, startDate, blackRule.CallCycle)
			if err != nil {
				return common.SystemInternalError
			}
			fqRequestCount = 1
			if frequencyCount+1 > int64(blackRule.CallCount) {
				fqHitCount = 1
				return common.OutOfFrequency
			}
		}
	}
	if blackRule.GatewayLevel != -1 {
		//TODO 根据Gateway 调用第三方接口
		sysGateway, err := model.GetSysGatewayImpl().GetByEnID(ctx, blackRule.GatewayLevel)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return common.NotFound
			}
			return common.SystemInternalError
		}
		gwRequestCount = 1
		isBlack, err := requestSysGateway(ctx, enID, sysGateway, callID, caller, callee)
		if err != nil {
			log.Errorf(ctx, "requestSysGateway error : callID: %s caller: %s callee: %s sysGateway : %+v err: %+v", callID, caller, callee, sysGateway, err)
		}
		if isBlack {
			gwHitCount = 1
			// 插入数据库
			err = model.GetMobileBlackApi().Insert(ctx, &model.MobileBlack{
				Mobile:    realCallee[3:],
				MobileAll: realCallee,
				MbLevel:   0,
				GwId:      sysGateway.NID,
				EnID:      enID,
			}, prefix)
			if err != nil {
				return common.SystemInternalError
			}
			return common.BlackMobile
		}
	}
	return common.StatusOK
}

func requestSysGateway(ctx context.Context, enID int, sg *model.SysGatewayInfo, callID interface{}, caller, callee string) (bool, error) {
	var isBlack bool
	var err error
	switch sg.GwType {
	case model.GwTypeVosHTTP:
		isBlack, err = vosHTTP(ctx, sg, enID, callID, caller, callee)
	case model.GwTypeVosRewrite:
		isBlack, err = vosRewrite(ctx, sg, enID, callID, caller, callee)
	case model.GwTypeDongyunBlack:
		isBlack, err = dongYun(ctx, sg, enID, callID, caller, callee)
	case model.GwTypeHuaxinVosBlack:
		isBlack, err = vosHuaXin(ctx, sg, callID, caller, callee)
	}
	return isBlack, err
}

func vosHuaXin(ctx context.Context, sg *model.SysGatewayInfo, callID interface{}, caller, callee string) (bool, error) {
	req := &http.HuaXinBlackCheck{
		CallID: callID,
		Caller: caller,
		Callee: callee,
	}
	res, err := req.Request(ctx, sg.GwUrl)
	if err != nil {
		return false, err
	}
	if res.ForbID == 1 && res.Status == 2001 {
		return true, nil
	}
	return false, nil
}

func dongYun(ctx context.Context, sg *model.SysGatewayInfo, enID int, callID interface{}, caller string, callee string) (bool, error) {
	str := getPreSign(sg.GwAk, callID, sg.GwPass)
	sign := utils.Encrypt(str)
	req := &http.SysGatewayDongYun{
		AK:     sg.GwAk,
		CallID: callID,
		Caller: caller,
		Callee: callee,
		Sign:   sign,
	}
	res, err := req.Request(ctx, sg.GwUrl)
	if err != nil {
		return false, err
	}
	if len(res.List) > 0 && res.List[0].Forbid == 1 {
		return true, nil
	}
	return false, nil
}

func getPreSign(ak string, callID interface{}, sk string) string {
	if id, ok := callID.(int64); ok {
		return fmt.Sprintf("%s%d%s", ak, id, sk)
	}
	if id, ok := callID.(string); ok {
		return fmt.Sprintf("%s%s%s", ak, id, sk)
	}
	if id, ok := callID.(float64); ok {
		return fmt.Sprintf("%s%.0f%s", ak, id, sk)
	}
	return ""
}

func vosHTTP(ctx context.Context, sg *model.SysGatewayInfo, enID int, callID interface{}, caller string, callee string) (bool, error) {
	req := &http.VOSHttpReq{
		CallID: callID,
		Caller: caller,
		Callee: callee,
	}
	res, err := req.Request(ctx, sg.GwUrl)
	if err != nil {
		return false, err
	}
	if res.ForbID == 1 {
		return true, nil
	}
	return false, nil
}

func vosRewrite(ctx context.Context, sg *model.SysGatewayInfo, enID int, callID interface{}, caller, callee string) (bool, error) {
	req := &http.HuaXinBlackRewrite{}
	req.RewriteE164Req.CallID = callID
	req.RewriteE164Req.CallerE164 = caller
	req.RewriteE164Req.CalleeE164 = callee
	res, err := req.Request(ctx, sg.GwUrl)
	if err != nil {
		return false, err
	}
	if res.RewriteE164Rsp.CalleeE164 != callee {
		return true, nil
	}
	return false, nil
}
