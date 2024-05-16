package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patnaikankit/Code-Cloud/server/pkg/controllers"
)

func Git(ctx *gin.RouterGroup) {
	router := ctx.Group("/git")

	router.POST("/clone", controllers.GitClone)
	router.DELETE("/delete", controllers.DeleteContainer)
}
