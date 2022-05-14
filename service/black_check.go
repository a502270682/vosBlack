package service

import (
	"context"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"vosBlack/common"
	"vosBlack/model"
	"vosBlack/utils"
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
	whiteNum, err := model.GetMobileWhitenumImpl().GetByWhiteNum(ctx, enID, realCallee)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// 在白名单内返回正常号码
	if whiteNum != nil {
		return true, nil
	}
	return false, nil
}

func CommonCheck(ctx context.Context, callee string, enID, ipID int) common.RespStatus {
	var prefix string
	if strings.HasPrefix(callee, "0") {
		prefix = "0"
	} else {
		prefix = callee[:3]
	}
	// 根据前缀和ip获取黑名单规则
	blackRule, err := model.GetEnterpriseBlacklistImpl().GetEnterpriseBlacklistByIPAndQianzhui(ctx, ipID, prefix)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.NotFound
		}
		return common.SystemInternalError
	}
	// 判断白名单
	if blackRule.IsWhitenum == 1 {
		isExist, err := IsWhiteNum(ctx, callee, enID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if isExist {
			return common.StatusOK
		}
	}
	// 判断呼叫时间段
	if blackRule.IsCalltime == 1 {
		callTimeList, err := model.GetEnterpriseCalltimelistImpl().GetByEnID(ctx, enID, blackRule.NID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if callTimeList != nil {
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
		mobilePatternList, err := model.GetMobilePatternImpl().GetListByMbLevel(ctx, blackRule.PatternLevel)
		if err != nil {
			return common.SystemInternalError
		}
		for _, value := range mobilePatternList {
			reg := regexp.MustCompile(value.Pattern)
			if reg.Match([]byte(callee)) {
				return common.PrettyNumber
			}
		}
	}
	// 判断本地黑名单
	if blackRule.BlacknumLevel != -1 {
		mobile, err := model.GetMobileBlackApi().GetOne(ctx, prefix, callee)
		if err != nil && err != gorm.ErrRecordNotFound {
			return common.SystemInternalError
		}
		if mobile != nil {
			return common.BlackMobile
		}
	}
	if blackRule.IsFrequency != 1 {
		//TODO 呼叫频率判断
	}
	if blackRule.GatewayLevel != -1 {
		//TODO 根据Gateway 调用第三方接口
	}
	return common.StatusOK
}
