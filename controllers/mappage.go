package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/maxiiot/vbaseBridge/config"
)

// VbaseMapPage return vbase map page url
// @summary vbase地图展示页URL
// @description vbase地图展示页URL
// @tags mappage
// @accept json
// @produce json
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /mappage [get]
func VbaseMapPage(c *gin.Context) {
	headers := make(url.Values)
	headers["appkey"] = []string{config.Cfg.VbaseServer.AppKey}
	headers["mapkey"] = []string{config.Cfg.VbaseServer.MapKey}
	pageurl := fmt.Sprintf("%s?%s", config.Cfg.VbaseServer.PageURL, headers.Encode())
	c.JSON(http.StatusOK, gin.H{
		"mappage": pageurl,
	})
}
