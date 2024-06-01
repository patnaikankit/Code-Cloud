package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/patnaikankit/Code-Cloud/server/pkg/routes"
	tools "github.com/patnaikankit/Code-Cloud/server/pkg/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// CORS configuration
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// WebSocket route
	router.GET("/ws", func(ctx *gin.Context) {
		done := tools.EstablishWS(ctx, &upgrader)
		if !done {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error establishing WebSocket connection",
			})
		}
	})

	// Default route for HTTP proxying
	router.NoRoute(handleHTTPProxy)

	// Initialize additional routes
	routes.Server(router)

	fmt.Println("Server started at port 5000")
	router.Run(":5000")
}

func handleHTTPProxy(ctx *gin.Context) {
	imageID := strings.Split(ctx.Request.Host, ".")[0]

	containerData, err := tools.ReadContainerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading container data",
			"error":   err.Error(),
		})
		return
	}

	containerInfo, exists := containerData[imageID]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Container Not Found",
		})
		return
	}

	serverURL := "http://localhost:" + strconv.Itoa(containerInfo.Port)
	internalServer, err := url.Parse(serverURL)
	if err != nil {
		fmt.Println("Error parsing internal server URL:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing internal server URL",
			"error":   err.Error(),
		})
		return
	}

	httpProxy := httputil.NewSingleHostReverseProxy(internalServer)
	httpProxy.ServeHTTP(ctx.Writer, ctx.Request)
}
