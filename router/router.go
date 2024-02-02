package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nullptr-z/gin-template/dao"
	"github.com/nullptr-z/gin-template/settings"
)

func Setup() *gin.Engine {
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery(), settings.LoggerFormateOutput)
	dao.InitializeDao()

	g.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Gin Template")
	})

	return g
}
