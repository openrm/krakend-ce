package krakend

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/proxy"
	router "github.com/luraproject/lura/router/gin"
)

func init() {
	router.RegisterRender("json_error", jsonErrorRender)
}

// Custom renderer to respond with the status code found in the resp metadata.
// See https://github.com/luraproject/lura/blob/master/router/gin/render.go for the defaults.
func jsonErrorRender(c *gin.Context, response *proxy.Response) {
	if response == nil {
		c.JSON(http.StatusOK, gin.H{})
	}
	status := http.StatusOK
	if response.Metadata.StatusCode > 0 {
		status = response.Metadata.StatusCode
	}
	c.JSON(status, response.Data)
}
