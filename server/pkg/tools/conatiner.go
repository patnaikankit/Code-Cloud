package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func WriteFileToContainer(containerID string, imageID string, filePath, data string) error {
	newFilePath := "./fileTemp" + imageID + filePath
	elements := strings.Split(newFilePath, "/")

	if len(elements) > 0 {
		lastIndex := len(elements) - 1
		elements = elements[:lastIndex]
	}

	dirPath := strings.Join(elements, "/")

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

	err = exec.Command("docker", "cp", newFilePath, containerID+":"+filePath).Run()

	newErr := os.Remove(newFilePath)
	if newErr != nil {
		return newErr
	}

	if err != nil {
		return err
	}

	return nil
}
