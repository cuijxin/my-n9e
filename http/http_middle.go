package http

import (
	"net/http"

	"github.com/cuijxin/my-n9e/pkg/ierr"
	"github.com/gin-gonic/gin"
)

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := loginUsername(c)
		c.Set("username", username)
		loginUser(c)
		c.Next()
	}
}

func admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := loginUsername(c)
		c.Set("username", username)

		user := loginUser(c)
		if user.Role != "Admin" {
			ierr.Bomb(http.StatusForbidden, "forbidden")
		}
		c.Next()
	}
}
