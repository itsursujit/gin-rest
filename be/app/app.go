package app

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loadRoutes(r *gin.Engine) {
	r.GET("/", HomePage)
}

func HomePage(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, "index", gin.H{
		"title": "Frontend title!",
	})
}
