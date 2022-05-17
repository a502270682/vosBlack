package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"vosBlack/common"
	"vosBlack/proto"
)

func getIP(ctx context.Context) string {
	return ctx.Value(common.IPCtxKey).(string)
}

func Error(c *gin.Context, respCode common.RespCode, respStatus common.RespStatus, inputType int, resp *proto.CommonReq) {
	switch inputType {
	case common.VOSRewrite:
		res := &proto.BlackScreeningRsp{
			Code:   respCode.Int(),
			Status: respStatus.Int(),
			CallID: resp.CallID,
			ForbID: resp.Callee,
		}
		c.JSON(http.StatusOK, res)
	case common.VOSHttp:
		res := &proto.BlackCheckRsp{}
		res.Code = respCode.Int()
		res.RewriteE164Rsp.Status = respStatus.Int()
		res.RewriteE164Rsp.CallID = resp.CallID
		res.RewriteE164Rsp.CallerE164 = resp.Caller
		res.RewriteE164Rsp.CalleeE164 = resp.Callee
		c.JSON(http.StatusOK, res)
	case common.SVOSHttp, common.DongyunHttp:
		res := &proto.BlackDongYunRsp{
			Code:   respCode.Int(),
			Status: respStatus.Int(),
			CallID: resp.CallID,
		}
		c.JSON(http.StatusOK, res)
	}
	return
}
