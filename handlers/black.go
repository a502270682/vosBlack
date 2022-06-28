package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"vosBlack/adapter/log"
	"vosBlack/common"
	"vosBlack/model"
	"vosBlack/proto"
	"vosBlack/service"
	"vosBlack/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

type CommonRsp struct {
	Code   int
	Msg    string
	Status int
}

func BlackDongYunHandler(c *gin.Context) {
	// 解析参数
	param, err := parseParam(common.DongyunHttp, c)
	if err != nil {
		if err == common.SignError {
			Error(c, common.SignErrorResp, common.NotFound, common.DongyunHttp, nil)
		} else if err == common.ReqParamError {
			Error(c, common.ParamError, common.NotFound, common.DongyunHttp, nil)
		} else if err == common.ReqAKError {
			Error(c, common.AKError, common.NotFound, common.DongyunHttp, nil)
		} else if err == common.ReqParamTypeError {
			Error(c, common.ParamTypeError, common.NotFound, common.DongyunHttp, nil)
		} else {
			Error(c, common.RespError, common.NotFound, common.DongyunHttp, nil)
		}
		// err = logic.UpsertEnterpriseApplyHourList(ctx, companyIPInfo.EnID, "", 1, 0, 0, 0, 0, 0, 0, 0, 0)
		// if err != nil {
		// 	log.Warnf(ctx, "UpsertEnterpriseApplyHourList fail, err:%+v", err)
		// }
		return
	}
	// 校验算法
	Check(c, param, common.DongyunHttp)
}

func BlackScreenHandler(c *gin.Context) {
	param, err := parseParam(common.VOSRewrite, c)
	if err != nil {
		if err == common.SignError {
			Error(c, common.SignErrorResp, common.NotFound, common.VOSRewrite, nil)
		} else if err == common.ReqParamError {
			Error(c, common.ParamError, common.NotFound, common.VOSRewrite, nil)
		} else if err == common.ReqAKError {
			Error(c, common.AKError, common.NotFound, common.VOSRewrite, nil)
		} else if err == common.ReqParamTypeError {
			Error(c, common.ParamTypeError, common.NotFound, common.VOSRewrite, nil)
		} else {
			Error(c, common.RespError, common.NotFound, common.VOSRewrite, nil)
		}
		// err = logic.UpsertEnterpriseApplyHourList(ctx, companyIPInfo.EnID, "", 1, 0, 0, 0, 0, 0, 0, 0, 0)
		// if err != nil {
		// 	log.Warnf(ctx, "UpsertEnterpriseApplyHourList fail, err:%+v", err)
		// }
		return
	}
	// 校验算法
	Check(c, param, common.VOSRewrite)
}

// 外部号码改写规则
func BlackCheckHandler(c *gin.Context) {
	param, err := parseParam(common.VOSHttp, c)
	if err != nil {
		if err == common.SignError {
			Error(c, common.SignErrorResp, common.NotFound, common.VOSHttp, nil)
		} else if err == common.ReqParamError {
			Error(c, common.ParamError, common.NotFound, common.VOSHttp, nil)
		} else if err == common.ReqAKError {
			Error(c, common.AKError, common.NotFound, common.VOSHttp, nil)
		} else if err == common.ReqParamTypeError {
			Error(c, common.ParamTypeError, common.NotFound, common.VOSHttp, nil)
		} else {
			Error(c, common.RespError, common.NotFound, common.VOSHttp, nil)
		}
		// err = logic.UpsertEnterpriseApplyHourList(ctx, companyIPInfo.EnID, "", 1, 0, 0, 0, 0, 0, 0, 0, 0)
		// if err != nil {
		// 	log.Warnf(ctx, "UpsertEnterpriseApplyHourList fail, err:%+v", err)
		// }
		return
	}
	// 校验算法
	Check(c, param, common.VOSHttp)
}

func haveBalance(ctx context.Context, enID int) bool {
	feel, err := service.GetEnterpriseFeelList(ctx, enID)
	if err != nil {
		return false
	}
	return feel.FeeIncome-feel.FeePayout > float64(0-feel.FeeCredit)
}

func Check(c *gin.Context, req *proto.CommonReq, inputType int) {
	ctx := c.Request.Context()
	ip := getClientIP(c)
	// 根据ip获取企业信息
	companyIPInfo, err := service.GetCompanyByIP(ctx, ip)
	// 检查企业实体是否存在和状态是否激活
	if err != nil || (companyIPInfo != nil && companyIPInfo.IStatus != model.IStatusActive) {
		c.JSON(http.StatusOK, CommonRsp{
			Code:   0,
			Status: common.NotFound.Int(),
		})
		return
	}
	enID := companyIPInfo.EnID
	ipID := companyIPInfo.NID
	// 判断余额
	if !haveBalance(ctx, enID) {
		Error(c, common.RespError, common.NoBalance, inputType, req)
		return
	}
	switch inputType {
	case common.VOSRewrite, common.VOSHttp:
		standloneCheck(c, req, inputType, ipID, enID)
	case common.SVOSHttp, common.DongyunHttp:
		loopCheck(c, req, ipID, enID)
	}
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
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				log.Errorf(context.Background(), "panic :%v", errors.WithStack(e))
			} else {
				log.Errorf(context.Background(), "panic :%v", err)
			}
			Error(c, common.RespError, common.NotFound, common.DongyunHttp, req)
			return
		}
	}()
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
			defer func() {
				if err := recover(); err != nil {
					if e, ok := err.(error); ok {
						log.Errorf(context.Background(), "panic :%v", errors.WithStack(e))
					} else {
						log.Errorf(context.Background(), "panic :%v", err)
					}
					syncList.AppendToArray(&proto.BlackDongYunDetail{
						Mobile: num,
						Msg:    "电话号异常",
						Status: -1,
					})
				}
			}()
			prefix, realCallee, phoneType := utils.GetPhone(num)
			respStatus := service.CommonCheck(ctx, prefix, realCallee, enID, ipID, req.CallID, req.Caller, num, phoneType)
			if respStatus == common.StatusOK {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 0,
					Msg:    common.CommonMobileType,
					Status: common.StatusOK,
				})
			} else if respStatus == common.OutOfFrequency {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 2,
					Msg:    common.OutMobileType,
					Status: common.OutOfFrequency,
				})
			} else {
				syncList.AppendToArray(&proto.BlackDongYunDetail{
					Mobile: realCallee,
					Forbid: 1,
					Msg:    common.SensitiveMobileType,
					Status: respStatus,
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
}

func standloneCheck(c *gin.Context, req *proto.CommonReq, inputType int, ipID int, enID int) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				log.Errorf(context.Background(), "panic :%v", errors.WithStack(e))
			} else {
				log.Errorf(context.Background(), "panic :%v", err)
			}
			Error(c, common.RespError, common.NotFound, inputType, req)
			return
		}
	}()
	ctx := c.Request.Context()
	prefix, realCallee, phoneType := utils.GetPhone(req.Callee)
	respStatus := service.CommonCheck(ctx, prefix, realCallee, enID, ipID, req.CallID, req.Caller, req.Callee, phoneType)
	if respStatus != common.StatusOK {
		Error(c, common.RespError, respStatus, inputType, req)
		return
	}
	if inputType == common.VOSRewrite {
		c.JSON(http.StatusOK, &proto.BlackScreeningRsp{
			Code:   1,
			CallID: req.CallID,
			Status: common.StatusOK.Int(),
		})
	} else {
		res := &proto.BlackCheckRsp{}
		res.Code = 1
		res.Memo = "success"
		res.Status = common.StatusOK.Int()
		res.RewriteE164Rsp.CallID = req.CallID
		res.RewriteE164Rsp.CallerE164 = req.Caller
		res.RewriteE164Rsp.CalleeE164 = req.Callee
		c.JSON(http.StatusOK, res)
	}
}

func parseParam(inputtype int, c *gin.Context) (*proto.CommonReq, error) {
	var req *proto.CommonReq
	var err error
	body := c.Request.Body
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, common.ReqParamError
	}
	log.Infof(c.Request.Context(), " req json :%v error: %v", string(reqBody), err)
	switch inputtype {
	case common.VOSRewrite:
		req, err = parseVOSRewriteParam(reqBody)
	case common.VOSHttp:
		req, err = parseVOSHttpParam(reqBody)
	case common.DongyunHttp, common.SVOSHttp:
		req, err = parseDongyunParam(reqBody)
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

func parseDongyunParam(reqBody []byte) (*proto.CommonReq, error) {
	req := proto.BlackDongYunReq{}
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, common.ReqParamError
	}
	if req.AK == "" {
		return nil, common.ReqParamTypeError
	}
	user, err := service.GetUserByAK(context.Background(), req.AK)
	if err != nil {
		return nil, err
	}
	if !checkSign(user.EnJkak, user.EnJkpass, req) {
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
	str := service.GetPreSign(userID, req.CallID, pass)
	return utils.Encrypt(str) == req.Sign
}

func parseVOSHttpParam(reqBody []byte) (*proto.CommonReq, error) {
	req := proto.BlackCheckReq{}
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, common.ReqParamError
	}
	if req.RewriteE164Req.CalleeE164 == "" {
		return nil, common.ReqParamTypeError
	}
	res := &proto.CommonReq{
		CallID: req.RewriteE164Req.CallID,
		Callee: req.RewriteE164Req.CalleeE164,
		Caller: req.RewriteE164Req.CallerE164,
	}
	return res, nil
}

func parseVOSRewriteParam(reqBody []byte) (*proto.CommonReq, error) {
	req := proto.BlackScreeningReq{}
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, common.ReqParamError
	}
	if req.Callee == "" {
		return nil, common.ReqParamTypeError
	}
	res := &proto.CommonReq{
		CallID: req.CallID,
		Callee: req.Callee,
		Caller: req.Caller,
	}
	return res, nil
}
