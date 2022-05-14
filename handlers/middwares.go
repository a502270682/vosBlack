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

func Error(c *gin.Context, respCode common.RespCode, respStatus common.RespStatus, inputType int) {
	switch inputType {
	case common.VOSRewrite:
		res := &proto.BlackScreeningRsp{
			Code:   respCode.Int(),
			Status: respStatus.Int(),
		}
		c.JSON(http.StatusOK, res)
	case common.VOSHttp:
		res := &proto.BlackCheckRsp{}
		res.Code = respCode.Int()
		res.RewriteE164Rsp.Status = respStatus.Int()
		c.JSON(http.StatusOK, res)
	case common.SVOSHttp, common.DongyunHttp:
		res := &proto.BlackDongYunRsp{
			Code:   respCode.Int(),
			Status: respStatus.Int(),
		}
		c.JSON(http.StatusOK, res)
	}
	return
}
