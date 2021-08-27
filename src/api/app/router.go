package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/config"
	"github.com/lmurature/melist-api/src/jobs"
	"github.com/onatm/clockwerk"
	"time"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
}

func StartApp() {
	mapUrls()

	c := clockwerk.New()
	c.Every(1 * time.Hour).Do(jobs.ItemsJobs)
	c.Start()

	router.Run(config.ApiPort)
}
