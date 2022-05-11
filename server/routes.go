package server

import (
	"net/http"
	"vosBlack/adapter/route"
	"vosBlack/handlers"

	"github.com/gin-gonic/gin"
)

func routes(engine *gin.Engine) {
	route.Route(engine, http.MethodGet, "/ping", handlers.PingHandler)
	route.Route(engine, http.MethodPost, "/black/check", handlers.BlackCheckHandler)
	route.Route(engine, http.MethodPost, "/black/screening", handlers.BlackScreeningHandler)
	route.Route(engine, http.MethodPost, "/black/dongyun", handlers.BlackDongYunHandler)
	//手机号黑名单管理功能
	admin := engine.Group("/admin")
	route.Route(admin, http.MethodPost, "/blackMobile/add", handlers.BlackMobileAddHandler)
	route.Route(admin, http.MethodPost, "/blackMobile/del", handlers.BlackMobileDelHandler)
	route.Route(admin, http.MethodPost, "/blackMobile/list", handlers.BlackMobileListHandler)

}
