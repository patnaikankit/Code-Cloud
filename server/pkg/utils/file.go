// to manipulate container data

package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/patnaikankit/Code-Cloud/server/pkg/models"
)

// to store container name and data
type ContainerMap map[string]models.ContainerInfo

// to parse the json data
func ReadContainerData() (ContainerMap, error) {
	absContextDirectory, err := filepath.Abs("../data/containers.info.json")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return nil, err
	}

	jsonFile, err := os.ReadFile(absContextDirectory)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	var containerMap ContainerMap

	err = json.Unmarshal(jsonFile, &containerMap)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}
	fmt.Println("JSON unmarshaled successfully")

	return containerMap, nil
}

// to write data to the json file
func WriteFile(data ContainerMap) error {
	absContextDirectory, err := filepath.Abs("../data/containers.info.json")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return err
	}

	jsonFile, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = os.WriteFile(absContextDirectory, jsonFile, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	fmt.Println("File written successfully to:", absContextDirectory)

	return nil
}

// retrieve container data
func FetchContainerData(imageName string, conatinerID string, port int) (string, error) {
	data, err := ReadContainerData()
	if err != nil {
		fmt.Println("Error reading container data:", err)
		return "", err
	}

	data[imageName] = models.ContainerInfo{
		ContainerID:   conatinerID,
		ContainerName: imageName,
		Port:          port,
	}
	fmt.Println("Container data updated:", data[imageName])

	err = WriteFile(data)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return "", err
	}
	fmt.Println("Container data written to file successfully")

	return conatinerID, nil
}

// to fetch avaliable ports
func FetchPort(startPort int) int {
	port := startPort

	for {
		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			// Windows command to check if the port is in use
			cmd = exec.Command("netstat", "-an")
		} else {
			// Unix command to check if the port is in use
			cmd = exec.Command("lsof", "-i", "tcp:"+strconv.Itoa(port))
		}

		output, err := cmd.Output()
		if err != nil {
			// Error running the command, assume port is available
			break
		}

		if runtime.GOOS == "windows" {
			// Check if the port appears in the netstat output
			if !strings.Contains(string(output), ":"+strconv.Itoa(port)) {
				break
			}
		} else {
			// If lsof returns output, port is in use
			if len(output) == 0 {
				break
			}
		}

		port++
	}

	return port
}
