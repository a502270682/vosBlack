package handlers

import (
	"context"
	"vosBlack/common"
)

func getIP(ctx context.Context) string {
	return ctx.Value(common.IPCtxKey).(string)
}
