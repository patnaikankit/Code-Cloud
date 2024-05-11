package tools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func EstablishWS(ctx *gin.Context, upgrader *websocket.Upgrader) bool {
	imageID := strings.Split(ctx.Request.Host, ".")[0]

	containerData, err := ReadContainersData()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error reading conatiner data!",
			"error":   err.Error(),
		})
		return false
	}

	containerInfo, temp := containerData[imageID]
	if !temp {
		ctx.JSON(404, gin.H{
			"message": "Conatiner Not Found!",
		})
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, http.Header{
		"Access-Control-Allow-Origin": []string{"*"},
	})

	if err != nil {
		fmt.Println("Error upgrading to websocket: ", err)
	}
	defer conn.Close()

	fmt.Println("Websocket Connection Established")

}
