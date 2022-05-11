package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vosBlack/adapter/route"
	"vosBlack/handlers"
)

func routes(engine *gin.Engine) {
	route.Route(engine, http.MethodGet, "/ping", handlers.PingHandler)
}
