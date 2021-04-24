package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/config"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
	router.Use(cors.Default())
}

func StartApp() {
	mapUrls()

	router.Run(config.ApiPort)
}
