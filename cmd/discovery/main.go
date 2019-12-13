package main

import (
	"fmt"
	"net/http"

	"github.com/aboglioli/big-brother/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Get()
	server := gin.Default()

	// Cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(cors.New(corsConfig))

	// Routes
	r := server.Group("/v1")
	{
		r.GET("/ping", ping)
	}

	server.Run(fmt.Sprintf(":%d", config.Discovery.Port))
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
