// deploying code from repo to container

package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	tools "github.com/patnaikankit/Code-Cloud/server/pkg/utils"
)

func GitClone(ctx *gin.Context) {
	// extract query parameters
	rootDir := ctx.Query("root-dir")
	stack := ctx.Query("stack")
	newRepo := uuid.New().String()

	// clone the repo in a temporary directory
	git.PlainClone("./tmp/"+newRepo, false, &git.CloneOptions{
		URL:      ctx.Query("git-link"),
		Progress: os.Stdout,
	})

	id, port, err := tools.CreateContainer(rootDir, stack, newRepo)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error creating container",
			"error":   err.Error(),
		})
		return
	}

	_, err = tools.FetchContainerData(newRepo, id, port)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Error creating container",
			"error":   err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"message": "Clone",
		"repo":    newRepo,
		"id":      id,
	})
}