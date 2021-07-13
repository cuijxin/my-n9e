package http

import (
	"github.com/cuijxin/my-n9e/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginPost(c *gin.Context) {
	var f loginForm
	bind(c, &f)

	user, err1 := models.PassLogin(f.Username, f.Password)
	if err1 == nil {
		if user.Status == 1 {
			renderMessage(c, "User disabled")
			return
		}
		session := sessions.Default(c)
		session.Set("username", f.Username)
		session.Save()
		renderData(c, user, nil)
		return
	}

	// password login fail, try ldap

	// password and ldap both fail
	renderMessage(c, err1)
}

func logoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("username", "")
	session.Save()
	renderMessage(c, nil)
}
