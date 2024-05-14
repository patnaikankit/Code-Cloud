// to manage container operations

package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// transfer data from host file system to the docker container
func WriteFileToContainer(containerID string, imageID string, filePath, data string) error {
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
