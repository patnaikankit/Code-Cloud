package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/patnaikankit/Code-Cloud/server/pkg/tools"
)

func main() {
	var timer time.Timer
	timer.C = make(<-chan time.Time)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// cors configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/*any", func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/ws" {
			done := tools.EstablishWS(ctx, &websocket.Upgrader{})
			if !done {
				ctx.JSON(500, gin.H{
					"message": "Error reading container data",
				})
			}
		}
	})

}
