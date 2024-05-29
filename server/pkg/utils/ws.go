// upgrade http to ws and manage communication with client
package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/patnaikankit/Code-Cloud/server/pkg/models"
)

func EstablishWS(ctx *gin.Context, upgrader *websocket.Upgrader) bool {
	// Retrieve container data
	imageID := strings.Split(ctx.Request.Host, ".")[0]

	containerData, err := ReadContainerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading container data!",
			"error":   err.Error(),
		})
		return false
	}

	containerInfo, exists := containerData[imageID]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Container Not Found!",
		})
		return false
	}

	// Upgrade the connection to WebSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return false
	}
	defer conn.Close()

	fmt.Println("WebSocket connection established")

	for {
		// Receive messages from the client
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		var command models.Command
		var output models.Output

		// Unmarshal the incoming command
		if err := json.Unmarshal(payload, &command); err != nil {
			output.Error = fmt.Sprintf("Error unmarshalling command: %v", err)
			if err := sendOutput(conn, messageType, output); err != nil {
				break
			}
			continue
		}

		// Execute command or write data to container
		if command.Data != "" {
			err := WriteToContainer(containerInfo.ContainerID, imageID, filepath.Join(command.Directory, command.IsFile), command.Data)
			if err != nil {
				output.Error = fmt.Sprintf("Error writing file to container: %v", err)
			}
		} else {
			cmd := exec.Command("docker", "exec", containerInfo.ContainerID, "sh", "-c", fmt.Sprintf("cd %s && %s", command.Directory, command.Command))
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				output.Error = fmt.Sprintf("Error starting command: %v", err)
			}
			if err := cmd.Wait(); err != nil {
				output.Error = fmt.Sprintf("Error waiting for command: %v", err)
			}
			if stdout.Len() > 0 {
				output.Output = stdout.String()
			}
			output.Type = command.Type

			if command.Type == "command" {
				pwd, err := getPwd(containerInfo.ContainerID, command.Directory, command.Command)
				if err != nil {
					output.Error = fmt.Sprintf("Error getting working directory: %v", err)
				}
				output.Directory = formatDirectory(pwd)
				output.OldDirectory = command.Directory
				output.Command = command.Command
			} else {
				output.Directory = command.Directory
			}
			output.IsFile = command.IsFile
		}

		// Send output back to client
		if err := sendOutput(conn, messageType, output); err != nil {
			break
		}
	}

	return true
}

// Helper function to format directory
func formatDirectory(pwd string) string {
	p := strings.Split(pwd, "/app")
	if len(p) > 1 {
		return "/app" + p[1]
	}
	return "/"
}

// Helper function to send output to the client
func sendOutput(conn *websocket.Conn, messageType int, output models.Output) error {
	outputData, err := json.Marshal(output)
	if err != nil {
		fmt.Println("Error marshalling output:", err)
		return err
	}
	if err := conn.WriteMessage(messageType, outputData); err != nil {
		fmt.Println("Error writing message:", err)
		return err
	}
	return nil
}

// Retrieve the directory path
func getPwd(containerID string, dir string, cmd string) (string, error) {
	pwd, err := exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("cd %s && %s && pwd", dir, cmd)).Output()
	return string(pwd), err
}
