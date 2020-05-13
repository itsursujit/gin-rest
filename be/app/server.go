package app

import (
	"fmt"
	"ginrest/be/app/Libraries"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

var (
	g errgroup.Group
)

func PrepareServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	initMiddleware(r)
	loadRoutes(r)
	loadView(r)
	return r
}

func Serve() {
	r := PrepareServer()
	serverInit(r)
}

func ServeViaProxy(startServers bool) {
	r := PrepareServer()
	serverProxyInit(r, startServers)
}

func ServeHost(host string, port int, r *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	return srv
}

func serverInit(r *gin.Engine) {
	host := envy.Get("APP_URL", "localhost")
	port, _ := strconv.Atoi(envy.Get("APP_PORT", "8080"))
	ServeHost(host, port, r)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func serverProxyInit(r *gin.Engine, startServers bool) {
	config := Libraries.ReadConfig("proxy")
	if startServers {
		for _, server := range config.Servers {
			ServeHost(server.Host, server.Port, r)
		}
	}
	Libraries.ServeWithProxy(config)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
