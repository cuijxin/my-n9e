package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cuijxin/my-n9e/models"
	"github.com/cuijxin/my-n9e/pkg/i18n"
	"github.com/cuijxin/my-n9e/pkg/ierr"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const defaultLimit = 20

func dangerous(v interface{}, code ...int) {
	ierr.Dangerous(v, code...)
}

func bomb(code int, format string, a ...interface{}) {
	ierr.Bomb(code, i18n.Sprintf(format, a...))
}

func bind(c *gin.Context, ptr interface{}) {
	dangerous(c.ShouldBindJSON(ptr), http.StatusBadRequest)
}

func urlParamStr(c *gin.Context, field string) string {
	val := c.Param(field)

	if val == "" {
		bomb(http.StatusBadRequest, "url param[%s] is blank", field)
	}

	return val
}

func urlParamInt64(c *gin.Context, field string) int64 {
	strval := urlParamStr(c, field)
	intval, err := strconv.ParseInt(strval, 10, 64)
	if err != nil {
		bomb(http.StatusBadRequest, "cannot convert %s to int64", strval)
	}
	return intval
}

func urlParamInt(c *gin.Context, field string) int {
	return int(urlParamInt64(c, field))
}

func queryStr(c *gin.Context, key string, defaultVal ...string) string {
	val := c.Query(key)
	if val != "" {
		return val
	}

	if len(defaultVal) == 0 {
		bomb(http.StatusBadRequest, "query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

func queryInt(c *gin.Context, key string, defaultVal ...int) int {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.Atoi(strv)
		if err != nil {
			bomb(http.StatusBadRequest, "cannot convert [%s] to int", strv)
		}
		return intv
	}

	if len(defaultVal) == 0 {
		bomb(http.StatusBadRequest, "query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

func offset(c *gin.Context, limit int) int {
	if limit <= 0 {
		limit = 10
	}

	page := queryInt(c, "p", 1)
	return (page - 1) * limit
}

func renderMessage(c *gin.Context, v interface{}, statusCode ...int) {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	if v == nil {
		c.JSON(code, gin.H{"err": ""})
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(code, gin.H{"err": i18n.Sprintf(t)})
	case error:
		c.JSON(code, gin.H{"err": t.Error()})
	}
}

func renderData(c *gin.Context, data interface{}, err error, statusCode ...int) {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	if err == nil {
		c.JSON(code, gin.H{"dat": data, "err": ""})
		return
	}

	renderMessage(c, err.Error(), code)
}

func cookieUsername(c *gin.Context) string {
	session := sessions.Default(c)

	value := session.Get("username")
	if value == nil {
		return ""
	}

	return value.(string)
}

func headerUsername(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}

	ut, err := models.UserTokenGet("token=?", strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return ""
	}

	if ut == nil {
		return ""
	}

	return ut.Username
}

func loginUsername(c *gin.Context) string {
	usernameInterface, has := c.Get("username")
	if has {
		return usernameInterface.(string)
	}

	username := cookieUsername(c)
	if username == "" {
		username = headerUsername(c)
	}

	if username == "" {
		ierr.Bomb(http.StatusUnauthorized, "unauthorized")
	}

	c.Set("username", username)
	return username
}

func loginUser(c *gin.Context) *models.User {
	username := loginUsername(c)

	user, err := models.UserGetByUsername(username)
	dangerous(err)

	if user == nil {
		ierr.Bomb(http.StatusUnauthorized, "unauthorized")
	}

	if user.Status == 1 {
		ierr.Bomb(http.StatusUnauthorized, "unauthorized")
	}

	return user
}

func User(id int64) *models.User {
	obj, err := models.UserGet("id=?", id)
	dangerous(err)

	if obj == nil {
		bomb(http.StatusNotFound, "No such user")
	}

	return obj
}
