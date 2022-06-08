package service

import (
	"context"
	"fmt"
	"strings"
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

func CommonCheck(ctx context.Context, prefix, realCallee string, enID, ipID int, callID, caller, callee string) common.RespStatus {
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
			err = logic.AddEnterpriseFqCache(ctx, enID, utils.GetLastNDay0TimeStamp(0), 1)
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
			curHour, curMintue := utils.GetCurHourAndMinute()
			if !(callTimeList.BeginHour <= curHour &&
				curHour <= callTimeList.EndHour &&
				callTimeList.BeginMinute <= curMintue &&
				curMintue < callTimeList.Edminute) {
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
		tablePrefix := realCallee[:3]
		if strings.HasPrefix(tablePrefix, "0") {
			tablePrefix = "0"
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
			frequencyCount, err := logic.GetEnterpriseFqFromStartDay(ctx, enID, startDate, blackRule.CallCycle)
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
		isBlack, err := requestSysGateway(ctx, enID, sysGateway.GwType, callID, caller, callee, sysGateway.GwAk, sysGateway.GwPass)
		if err != nil {
			return common.SystemInternalError
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

func requestSysGateway(ctx context.Context, enID int, gwType model.GwType, callID, caller, callee string, ak, pass string) (bool, error) {
	var isBlack bool
	var err error
	switch gwType {
	case 1:
		isBlack, err = vosRewrite(ctx, enID, callID, caller, callee)
	case 2:
		isBlack, err = vosYunXuntong(ctx, enID, callID, caller, callee)
	case 3:
		isBlack, err = dongYun(ctx, enID, callID, caller, callee, ak, pass)
	}
	return isBlack, err
}

func dongYun(ctx context.Context, enID int, callID string, caller string, callee, ak, pass string) (bool, error) {
	str := strings.ToLower(fmt.Sprintf("%s%s%s", ak, callID, pass))
	sign := utils.Encrypt(str)
	req := &http.SysGatewayDongYun{
		AK:     ak,
		CallID: callID,
		Caller: caller,
		Callee: callee,
		Sign:   sign,
	}
	res, err := req.Request(ctx)
	if err != nil {
		return false, err
	}
	if res.List[0].Forbid == 0 {
		return true, nil
	}
	return false, nil
}

func vosYunXuntong(ctx context.Context, enID int, callID string, caller string, callee string) (bool, error) {
	req := &http.HuaXinBlackCheck{
		CallID: callID,
		Caller: caller,
		Callee: callee,
	}
	res, err := req.Request(ctx)
	if err != nil {
		return false, err
	}
	if res.Status != 2000 {
		return true, nil
	}
	return false, nil
}

func vosRewrite(ctx context.Context, enID int, callID, caller, callee string) (bool, error) {
	req := &http.HuaXinBlackRewrite{}
	req.RewriteE164Req.CallID = callID
	req.RewriteE164Req.CallerE164 = caller
	req.RewriteE164Req.CalleeE164 = callee
	res, err := req.Request(ctx)
	if err != nil {
		return false, err
	}
	if res.RewriteE164Rsp.Status != 2000 {
		return true, nil
	}
	return false, nil
}
