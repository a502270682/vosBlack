package server

import (
	"net/http"
	"vosBlack/adapter/route"
	"vosBlack/handlers"

	"github.com/gin-gonic/gin"
)

func routes(engine *gin.Engine) {
	black := engine.Group("/black")
	black.POST("/check", handlers.BlackCheckHandler)
	black.POST("/screening", handlers.BlackCheckHandler)
	black.POST("/dongyun", handlers.BlackCheckHandler)
	//手机号黑名单管理功能
	admin := engine.Group("/admin")
	route.Route(admin, http.MethodPost, "/blackMobile/add", handlers.BlackMobileAddHandler)
	route.Route(admin, http.MethodPost, "/blackMobile/del", handlers.BlackMobileDelHandler)
	route.Route(admin, http.MethodGet, "/blackMobile/list", handlers.BlackMobileListHandler)
	route.Route(admin, http.MethodGet, "/blackMobile/info", handlers.BlackMobileInfoHandler)
}
