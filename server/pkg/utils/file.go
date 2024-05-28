// to manipulate container data

package tools

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/patnaikankit/Code-Cloud/server/pkg/models"
)

// to store container name and data
type ContainerMap map[string]models.ContainerInfo

// to parse the json data
func ReadContainerData() (ContainerMap, error) {
	absContextDirectory, err := filepath.Abs("../data/containers.info.json")
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(absContextDirectory)
	if err != nil {
		return nil, err
	}

	var conatinerMap ContainerMap

	err = json.Unmarshal(jsonFile, &conatinerMap)
	if err != nil {
		return nil, err
	}

	return conatinerMap, nil
}

// to write data to the json file
func WriteFile(data ContainerMap) error {
	absContextDirectory, err := filepath.Abs("../data/containers.info.json")
	if err != nil {
		return err
	}

	jsonFile, err := json.Marshal(absContextDirectory)
	if err != nil {
		return err
	}

	err = os.WriteFile(absContextDirectory, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// retrieve container data
func FetchContainerData(imageName string, conatinerID string, port int) (string, error) {
	data, err := ReadContainerData()
	if err != nil {
		return "", err
	}

	data[imageName] = models.ContainerInfo{
		ContainerID:   conatinerID,
		ContainerName: imageName,
		Port:          port,
	}

	err = WriteFile(data)

	if err != nil {
		return "", err
	}

	return conatinerID, nil
}

// to fetch avaliable ports
func FetchPort(port int) int {
	for {
		port++
		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			// Windows command to check if the port is in use
			cmd = exec.Command("netstat", "-an", "|", "findstr", ":"+strconv.Itoa(port))
		} else {
			// Unix for linux or mac command to check if the port is in use
			cmd = exec.Command("lsof", "-i", "tcp:"+strconv.Itoa(port))
		}

		err := cmd.Run()
		if err != nil {
			break
		}
	}

	return port
}
