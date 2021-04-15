package ping

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	pong string = "pong"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, pong)
}
