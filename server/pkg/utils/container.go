// to manage container operations

package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// create and start the container
func CreateContainer(rootDir string, stack string, imageName string) (string, int, error) {
	// location where the dockerfile will be copied
	code := ""
	hostPort := FetchPort(8080)

	if rootDir == "" {
		code = "./tmp/" + imageName
	} else {
		code = "./tmp/" + imageName + "/" + rootDir
	}

	image := imageName + ":latest"

	err := exec.Command("cp", "./pkg/dockerFiles/"+stack+"/Dockerfile", code).Run()
	if err != nil {
		return "", 0, err
	}

	// build image
	err = exec.Command("docker", "build", "-t", image, code).Run()
	if err != nil {
		fmt.Println("building", err)
		return "", 0, err
	}

	// build container
	createCmd := exec.Command("docker", "create", "--name", imageName, "-p", strconv.Itoa(hostPort)+":3000", image)
	output, err := createCmd.CombinedOutput()
	if err != nil {
		return "", 0, err
	}

	conatainerID := string(bytes.TrimSpace(output))

	startCmd := exec.Command("docker", "start", conatainerID)
	_, err = startCmd.CombinedOutput()
	if err != nil {
		return "", 0, err
	}

	return conatainerID, hostPort, nil
}

// transfer data from host file system to the docker container
func WriteToContainer(containerID string, imageID string, filePath, data string) error {
	// temporarily storing the file in the host's file system
	newFilePath := "./fileTemp" + imageID + filePath
	elements := strings.Split(newFilePath, "/")

	if len(elements) > 0 {
		lastIndex := len(elements) - 1
		elements = elements[:lastIndex]
	}

	dirPath := strings.Join(elements, "/")

	// make new directory if it doesn't exist
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("Error creating directory -> ", err)
		return err
	}

	err = os.WriteFile(newFilePath, []byte(data), 6044)
	if err != nil {
		fmt.Println("Error writing file ->", err)
		return err
	}

	// copy files from host to container
	err = exec.Command("docker", "cp", newFilePath, containerID+":"+filePath).Run()

	// remove the temporary files from host
	newErr := os.Remove(newFilePath)
	if newErr != nil {
		return newErr
	}

	if err != nil {
		return err
	}

	return nil
}

// delete the conatiner and the associated docker image
func DeleteImageAndContainer(imageID string, containerID string) error {
	err := exec.Command("docker", "stop", containerID).Run()
	if err != nil {
		return err
	}

	err = exec.Command("docker", "rm", containerID).Run()
	if err != nil {
		return err
	}

	err = exec.Command("docker", "rmi", imageID).Run()
	if err != nil {
		return err
	}

	return nil
}
