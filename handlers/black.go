package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"vosBlack/adapter/log"
	"vosBlack/adapter/logic"
	"vosBlack/common"
	"vosBlack/proto"
	"vosBlack/service"
	"vosBlack/utils"
)

// 外部号码改写规则
func BlackCheckHandler(c *gin.Context) {
	ctx := c.Request.Context()
	ip := c.ClientIP()
	// 根据ip获取企业信息
	companyIPInfo, err := service.GetCompanyByIP(ctx, ip)
	if err != nil {
		c.JSON(http.StatusOK, struct {
			Code   int
			Status int
		}{
			Code:   0,
			Status: common.NotFound.Int(),
		})
		return
	}
	// 判断余额
	if !haveBalance(ctx, companyIPInfo.EnID) {
		Error(c, common.RespError, common.NoBalance, companyIPInfo.Inputtype)
		return
	}
	// 解析参数
	param, err := parseParam(companyIPInfo.Inputtype, c)
	if err != nil {
		err = logic.UpsertEnterpriseApplyHourList(ctx, companyIPInfo.EnID, "", 1, 0, 0, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Warnf(ctx, "UpsertEnterpriseApplyHourList fail, err:%+v", err)
		}
		if err == common.SignError {
			Error(c, common.SignErrorResp, common.NotFound, companyIPInfo.Inputtype)
		} else if err == common.ReqParamError {
			Error(c, common.ParamError, common.NotFound, companyIPInfo.Inputtype)
		} else if err == common.ReqAKError {
			Error(c, common.AKError, common.NotFound, companyIPInfo.Inputtype)
		} else {
			Error(c, common.RespError, common.NotFound, companyIPInfo.Inputtype)
		}
		return
	}
	// 校验算法
	Check(c, param, companyIPInfo.Inputtype, companyIPInfo.NID, companyIPInfo.EnID)
	return
}

func haveBalance(ctx context.Context, enID int) bool {
	feel, err := service.GetEnterpriseFeelList(ctx, enID)
	if err != nil {
		return false
	}
	return feel.FeeIncome-feel.FeePayout >= float64(0-feel.FeeCredit)
}

func Check(c *gin.Context, req *proto.CommonReq, inputType, ipID, enID int) {
	switch inputType {
	case common.VOSRewrite, common.VOSHttp:
		standloneCheck(c, req, inputType, ipID, enID)
	case common.SVOSHttp, common.DongyunHttp:
		loopCheck(c, req, ipID, enID)
	}
	return
}

type SyncList struct {
	List []*proto.BlackDongYunDetail
	sync.RWMutex
}

func (sl *SyncList) AppendToArray(detail *proto.BlackDongYunDetail) {
	sl.Lock()
	defer sl.Unlock()
	sl.List = append(sl.List, detail)
}

func loopCheck(c *gin.Context, req *proto.CommonReq, ipID int, enID int) {
	ctx := c.Request.Context()
	calleeArr := strings.Split(req.Callee, ",")
	wg := sync.WaitGroup{}
	syncList := &SyncList{
		List: []*proto.BlackDongYunDetail{},
	}
	for _, callee := range calleeArr {
		wg.Add(1)
		go func(num string) {
			defer wg.Done()
			matchs := utils.FindStringSubmatch(num)
			// 获取真实的手机号
			realCallee := matchs[5]
			respStatus := service.CommonCheck(ctx, realCallee, enID, ipID, req.CallID, req.Caller, req.Callee)
			if respStatus == common.StatusOK {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 0,
					Msg:    common.CommonMobileType,
				})
			} else if respStatus == common.OutOfFrequency {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 2,
					Msg:    common.OutMobileType,
				})
			} else {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 1,
					Msg:    common.SensitiveMobileType,
				})
			}

		}(callee)
	}
	wg.Wait()
	res := &proto.BlackDongYunRsp{
		Code:          1,
		Msg:           "success",
		CallID:        req.CallID,
		TransactionID: req.Caller,
		List:          syncList.List,
	}
	c.JSON(http.StatusOK, res)
	return
}

func standloneCheck(c *gin.Context, req *proto.CommonReq, inputType int, ipID int, enID int) {
	ctx := c.Request.Context()
	respStatus := service.CommonCheck(ctx, req.Callee, enID, ipID, req.CallID, req.Caller, req.Callee)
	if respStatus != common.StatusOK {
		Error(c, common.RespError, respStatus, inputType)
		return
	}
	if inputType == common.VOSRewrite {
		c.JSON(http.StatusOK, &proto.BlackScreeningRsp{
			Code:   1,
			CallID: req.CallID,
			ForbID: req.Callee,
			Status: common.StatusOK.Int(),
		})
	} else {
		res := &proto.BlackCheckRsp{}
		res.Code = 1
		res.Memo = "success"
		res.RewriteE164Rsp.Status = common.StatusOK.Int()
		res.RewriteE164Rsp.CallID = req.CallID
		res.RewriteE164Rsp.CallerE164 = req.Caller
		res.RewriteE164Rsp.CalleeE164 = req.Callee
		c.JSON(http.StatusOK, res)
	}
	return
}

func parseParam(inputtype int, c *gin.Context) (*proto.CommonReq, error) {
	var req *proto.CommonReq
	var err error
	switch inputtype {
	case common.VOSRewrite:
		req, err = parseVOSRewriteParam(c)
	case common.VOSHttp:
		req, err = parseVOSHttpParam(c)
	case common.DongyunHttp, common.SVOSHttp:
		req, err = parseDongyunParam(c)
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

func parseDongyunParam(c *gin.Context) (*proto.CommonReq, error) {
	ctx := c.Request.Context()
	req := proto.BlackDongYunReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return nil, err
	}
	user, err := service.GetUserByUserID(ctx, req.AK)
	if err != nil {
		return nil, err
	}
	if !checkSign(user.UserID, user.UserPass, req) {
		return nil, common.SignError
	}
	res := &proto.CommonReq{
		CallID: req.CallID,
		Callee: req.Callee,
		Caller: req.Caller,
	}
	return res, nil
}

func checkSign(userID string, pass string, req proto.BlackDongYunReq) bool {
	str := fmt.Sprintf("%s%s%s", userID, req.CallID, pass)
	return utils.Encrypt(str) == req.Sign
}

func parseVOSHttpParam(c *gin.Context) (*proto.CommonReq, error) {
	req := proto.BlackCheckReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return nil, common.ReqParamError
	}
	res := &proto.CommonReq{
		CallID: strconv.Itoa(req.RewriteE164Req.CallID),
		Callee: req.RewriteE164Req.CalleeE164,
		Caller: req.RewriteE164Req.CallerE164,
	}
	return res, nil
}

func parseVOSRewriteParam(c *gin.Context) (*proto.CommonReq, error) {
	req := proto.BlackScreeningReq{}
	err := c.BindJSON(&req)
	if err != nil {
		return nil, common.ReqParamError
	}
	res := &proto.CommonReq{
		CallID: req.CallID,
		Callee: req.Callee,
		Caller: req.Caller,
	}
	return res, nil
}
