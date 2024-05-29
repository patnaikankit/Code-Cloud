// to manage container operations

package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/otiai10/copy"
)

// create and start the container
func CreateContainer(rootDir, stack, imageName string) (string, int, error) {
	code := ""
	hostPort := FetchPort(8080)

	if rootDir == "" {
		code = "./tmp/" + imageName
	} else {
		code = "./tmp/" + imageName + "/" + rootDir
	}

	image := imageName + ":latest"

	fmt.Println(image)

	fmt.Println("--------------------------------------------------------------------------------------")

	dockerfilePath := "../pkg/docker/" + stack + "/Dockerfile"
	destinationPath := filepath.Join(code, "Dockerfile")

	err := copy.Copy(dockerfilePath, destinationPath)
	if err != nil {
		fmt.Println("Error copying Dockerfile:", err)
		return "", 0, err
	}

	// Build the Docker image
	buildCmd := exec.Command("docker", "build", "-t", image, code)
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error building Docker image:", string(buildOutput), err)
		return "", 0, err
	}
	fmt.Println("Build output:", string(buildOutput))

	// Create the Docker container
	createCmd := exec.Command("docker", "create", "--name", imageName, "-p", strconv.Itoa(hostPort)+":3000", image)
	createOutput, err := createCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error creating Docker container:", string(createOutput), err)
		return "", 0, err
	}
	fmt.Println("Create output:", string(createOutput))

	containerID := string(bytes.TrimSpace(createOutput))

	// Start the Docker container
	startCmd := exec.Command("docker", "start", containerID)
	startOutput, err := startCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error starting Docker container:", string(startOutput), err)
		return "", 0, err
	}
	fmt.Println("Start output:", string(startOutput))

	return containerID, hostPort, nil
}

// transfer data from host file system to the docker container
func WriteToContainer(containerID string, imageID string, filePath, data string) error {
	// Construct temporary file path in the host's file system.
	newFilePath := filepath.Join(".", "fileTmp"+imageID, filePath)
	dirPath := filepath.Dir(newFilePath)

	// Create directory if it doesn't exist.
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("Error creating directory -> ", err)
		return err
	}

	// Write data to the temporary file.
	err = os.WriteFile(newFilePath, []byte(data), 0644)
	if err != nil {
		fmt.Println("Error writing file ->", err)
		return err
	}

	// Create a temporary directory for copying.
	tempDir := filepath.Join(".", "fileTmp"+imageID)
	defer os.RemoveAll(tempDir)

	// Copy the temporary file to the container.
	hostPath := filepath.Join(tempDir, filePath)
	err = copy.Copy(hostPath, newFilePath)
	if err != nil {
		fmt.Println("Error copying file ->", err)
		return err
	}

	err = exec.Command("docker", "cp", newFilePath, containerID+":"+filePath).Run()
	if err != nil {
		fmt.Println("Error copying file to container ->", err)
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
