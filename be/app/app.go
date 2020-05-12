package app

import (
	"context"
	"fmt"
	"ginrest/be/app/Http/Middleware"
	"ginrest/be/app/Libraries"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)
var log *logrus.Logger
func Serve() {
	r := gin.New()
	loadView(r)
	initMiddleware(r)
	loadRoutes(r)
	serverInit(r)
}

func serverInit(r *gin.Engine) {
	host := envy.Get("APP_URL", "localhost")
	port := envy.Get("APP_PORT", "8080")
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //nolint:staticcheck
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select { //nolint:gosimple
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}
	log.Println( "Server exiting")
}

func initMiddleware(r *gin.Engine) {
	log = Libraries.NewLogger()

	log.SetFormatter(&logrus.JSONFormatter{})

	r.Use(Middleware.RequestId(), Middleware.HttpLogger(log), gin.Recovery())
}

func loadRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	r.POST("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})
	r.GET("/page", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
	})
}

func loadView(r *gin.Engine) {
	r.HTMLRender = ginview.New(goview.Config{
		Root:      "./fe/views",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format(strconv.Itoa(time.Now().Year()))
			},
		},
		DisableCache: false,
	})
}