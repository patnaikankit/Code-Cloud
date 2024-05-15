// to update code files in the container after changes
package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/patnaikankit/Code-Cloud/server/pkg/models"
	tools "github.com/patnaikankit/Code-Cloud/server/pkg/utils"
)

func CodeUpdate(ctx *gin.Context) {
	imageID := strings.Split(ctx.Request.Host, ".")[0]

	var files []models.File
	err := ctx.ShouldBindJSON(&files)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	containerData, err := tools.ReadContainerData()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error reading container data",
			"error":   err.Error(),
		})
		return
	}

	containerInfo, temp := containerData[imageID]
	if !temp {
		ctx.JSON(404, gin.H{
			"message": "Container Not Found!",
		})
		return
	}

	// commit changes in the container
	for _, file := range files {
		err = tools.WriteToContainer(containerInfo.ContainerID, imageID, file.FilePath, file.Data)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error writing file to container",
				"error":   err.Error(),
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message": "Updated",
	})
}
