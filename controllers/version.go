package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var version string

func SetVersion(v string) {
	version = v
}

// @summary version
// @tags version
// @accept json
// @produce json
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /version [get]
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": version,
	})
}
