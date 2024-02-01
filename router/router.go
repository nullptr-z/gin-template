package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

	g.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Gin Template")
	})

	return g
}
