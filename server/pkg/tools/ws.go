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
)

type Command struct {
	Directory string `json:"directory"`
	Command   string `json:"command"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	IsFile    string `json:"isFile"`
}

type Output struct {
	OldDirectory string `json:"oldDirectory"`
	Directory    string `json:"directory"`
	Output       string `json:"ouput"`
	Error        string `json:"error"`
	Type         string `json:"typ"`
	IsFile       string `json:"isFile"`
	Command      string `json:"command"`
}

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

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return false
		}

		var command Command
		var output Output

		execCommand := string(p)
		json.Unmarshal([]byte(execCommand), &command)
		var cmd *exec.Cmd
		if command.Data != "" {
			err := WriteFileToContainer(containerInfo.ContainerID, imageID, filepath.Join(command.Directory, command.IsFile), command.Data)

			if err != nil {
				output.Error = fmt.Sprintf("Error writing file to conatiner: %v\n", err)
			}
		} else {
			cmd = exec.Command("docker", "exec", containerInfo.ContainerID, "sh", "-c", "cd"+command.Directory+"&& "+command.Command)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				output.Error = fmt.Sprintf("Error starting conatiner: %v\n", err)
			}
			if err := cmd.Wait(); err != nil {
				output.Error = fmt.Sprintf("Error waiting for conatiner: %v\n", err)
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
				return false
			}
		}
	}
}

func getPwd(containerID string, dir string, cmd string) (string, error) {
	pwd, err := exec.Command("docker", "exec", containerID, "sh", "-c", "cd"+dir+"&& "+cmd+" && pwd").Output()
	return string(pwd), err
}
