package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index return basicauth index
func Index(c *gin.Context) {
	b, err := ioutil.ReadFile("./ui/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err)
		return
	}

	c.Writer.Write(b)
}
