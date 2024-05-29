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

// to send data and execute commands in container
func EstablishWS(ctx *gin.Context, upgrader *websocket.Upgrader) bool {
	// retrieve container data
	imageID := strings.Split(ctx.Request.Host, ".")[0]

	containerData, err := ReadContainerData()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error reading container data!",
			"error":   err.Error(),
		})
		return false
	}

	containerInfo, temp := containerData[imageID]
	if !temp {
		ctx.JSON(404, gin.H{
			"message": "Container Not Found!",
		})
		return false
	}

	// upgrading the connection to websocket from http
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, http.Header{
		"Access-Control-Allow-Origin": []string{"*"},
	})

	if err != nil {
		fmt.Println("Error upgrading to websocket: ", err)
		return false
	}
	defer conn.Close()

	fmt.Println("Websocket Connection Established")

	for {
		// receiving signals from client
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return false
		}

		var command models.Command
		var output models.Output

		execCommand := string(payload)
		json.Unmarshal([]byte(execCommand), &command)
		var cmd *exec.Cmd
		// write operation
		if command.Data != "" {
			err := WriteToContainer(containerInfo.ContainerID, imageID, filepath.Join(command.Directory, command.IsFile), command.Data)

			if err != nil {
				output.Error = fmt.Sprintf("Error writing file to container: %v\n", err)
			}
		} else {
			// command execution inside container
			cmd = exec.Command("docker", "exec", containerInfo.ContainerID, "sh", "-c", "cd "+command.Directory+" && "+command.Command)

			// to store the output of the executed command
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				output.Error = fmt.Sprintf("Error starting container: %v\n", err)
			}
			if err := cmd.Wait(); err != nil {
				output.Error = fmt.Sprintf("Error waiting for container: %v\n", err)
			}
			if stdout.Len() > 0 {
				output.Output = stdout.String()
			}
			output.Type = command.Type
			if command.Type == "command" {
				pwd, err := getPwd(
					containerInfo.ContainerID,
					command.Directory,
					command.Command,
				)
				if err != nil {
					output.Error = fmt.Sprintf("Error getting pwd: %v \n", err)
				}

				p := strings.Split(pwd, "/app")
				if len(p) > 1 {
					output.Directory = "/app" + p[1]
				} else {
					output.Directory = "/"
				}
				output.OldDirectory = command.Directory
				output.Command = command.Command
			} else {
				output.Directory = command.Directory
			}
			output.IsFile = command.IsFile
			st, _ := json.Marshal(output)
			if err := conn.WriteMessage(messageType, st); err != nil {
				conn.Close()
				return false
			}
		}
	}
}

// retrieve the directory path
func getPwd(containerID string, dir string, cmd string) (string, error) {
	pwd, err := exec.Command("docker", "exec", containerID, "sh", "-c", "cd "+dir+" && "+cmd+" && pwd").Output()
	return string(pwd), err
}
