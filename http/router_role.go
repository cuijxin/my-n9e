package http

import (
	"github.com/cuijxin/my-n9e/models"
	"github.com/gin-gonic/gin"
)

func rolesGet(c *gin.Context) {
	lst, err := models.RoleGetsAll()
	renderData(c, lst, err)
}
