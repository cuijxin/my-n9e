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
		guest.GET("/auth/logout", logoutGet)
	}

	// for brower, expose location in nginx.conf
	pages := r.Group("/api/n9e", csrfMid)
	{
		pages.GET("/csrf", func(c *gin.Context) {
			renderData(c, csrf.GetToken(c), nil)
		})

		pages.GET("/roles", rolesGet)
		pages.GET("/self/profile", selfProfileGet)
		pages.PUT("/self/profile", selfProfilePut)
		pages.PUT("/self/password", selfPasswordPut)
		pages.GET("/self/token", selfTokenGets)
		pages.POST("/self/token", selfTokenPost)
		pages.PUT("/self/token", selfTokenPut)

		pages.GET("/users", login(), userGets)
		pages.POST("/users", admin(), userAddPost)
		pages.GET("/user/:id/profile", login(), userProfileGet)
		pages.PUT("/user/:id/profile", admin(), userProfilePut)
		pages.PUT("/user/:id/status", admin(), userStatusPut)
		pages.PUT("/user/:id/password", admin(), userPasswordPut)
	}
}
