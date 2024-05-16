// to delete container and image and clear the directory associated with it

package controllers

import (
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	tools "github.com/patnaikankit/Code-Cloud/server/pkg/utils"
)

func DeleteContainer(ctx *gin.Context) {
	imageID := ctx.Query("image-id")

	containerData, err := tools.ReadContainerData()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error reading container's data",
			"error":   err.Error(),
		})
		return
	}

	containerInfo, temp := containerData[imageID]
	if !temp {
		ctx.JSON(404, gin.H{
			"message": "Container Not Found",
		})
	}

	err = tools.DeleteImageAndContainer(imageID+":latest", containerInfo.ContainerID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error deleting container",
			"error":   err.Error(),
		})
		return
	}

	delete(containerData, imageID)
	err = tools.WriteFile(containerData)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error writing container data",
			"err":     err.Error(),
		})
		return
	}

	// deleting directory
	absPath, err := filepath.Abs("./tmp/" + imageID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error getting absolute path",
			"error":   err.Error(),
		})
		return
	}

	err = exec.Command("rm", "-r", absPath).Run()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error deleting directory",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Deleted",
	})
}
