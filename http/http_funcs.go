package http

import (
	"net/http"

	"github.com/cuijxin/my-n9e/pkg/i18n"
	"github.com/cuijxin/my-n9e/pkg/ierr"
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
