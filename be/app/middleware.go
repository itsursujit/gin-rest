package app

import (
	"ginrest/be/app/Http/Middleware"
	"ginrest/be/app/Libraries"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func initMiddleware(r *gin.Engine) {
	log = Libraries.NewLogger()
	log.SetFormatter(&logrus.JSONFormatter{})
	r.Use(Middleware.RequestId(), Middleware.HttpLogger(log), gin.Recovery())
	r.Use(static.Serve("/assets", static.LocalFile("./assets", false)))
}
