package handlers

import (
	"context"
	"fmt"
	"net/http"
	"vosBlack/common"
	"vosBlack/proto"

	"github.com/gin-gonic/gin"
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
		}
		if resp != nil {
			res.CallID = resp.CallID
			res.ForbID = 1
		}
		c.JSON(http.StatusOK, res)
	case common.VOSHttp:
		res := &proto.BlackCheckRsp{}
		res.Code = respCode.Int()
		res.Status = respStatus.Int()
		if resp != nil {
			res.RewriteE164Rsp.CallID = resp.CallID
			res.RewriteE164Rsp.CallerE164 = resp.Caller
			res.RewriteE164Rsp.CalleeE164 = fmt.Sprintf("%d-%s", respStatus.Int(), resp.Callee)
		}
		c.JSON(http.StatusOK, res)
	case common.SVOSHttp, common.DongyunHttp:
		res := &proto.BlackDongYunRsp{
			Code:   respCode.Int(),
			Status: respStatus.Int(),
		}
		if resp != nil {
			res.CallID = resp.CallID
		}
		c.JSON(http.StatusOK, res)
	}
	return
}
