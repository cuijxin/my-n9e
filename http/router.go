package http

import (
	"fmt"
	"os"

	"github.com/cuijxin/my-n9e/config"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/csrf"
)

func configRoutes(r *gin.Engine) {
	csrfMid := csrf.Middleware(csrf.Options{
		Secret: config.Config.HTTP.CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.JSON(452, gin.H{"err": "csrf token mismatch"})
			c.Abort()
		},
	})

	if config.Config.HTTP.Pprof {
		pprof.Register(r, "/api/debug/pprof")
	}

	guest := r.Group("api/n9e")
	{
		guest.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		guest.GET("/pid", func(c *gin.Context) {
			c.String(200, fmt.Sprintf("%d", os.Getpid()))
		})
		guest.GET("/addr", func(c *gin.Context) {
			c.String(200, c.Request.RemoteAddr)
		})
		guest.POST("/auth/login", loginPost)
	}

	// for brower, expose location in nginx.conf
	pages := r.Group("/api/n9e", csrfMid)
	{
		pages.GET("/csrf", func(c *gin.Context) {
			renderData(c, csrf.GetToken(c), nil)
		})
	}
}
