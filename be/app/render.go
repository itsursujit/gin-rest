package app

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/foolin/goview/supports/gorice"
	"github.com/gin-gonic/gin"
)

func loadView(r *gin.Engine) {
	basic := gorice.NewWithConfig(rice.MustFindBox("../../fe/views"), goview.Config{
		Root:         "./fe/views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		DisableCache: false,
	})
	r.HTMLRender = ginview.Wrap(basic)
}
