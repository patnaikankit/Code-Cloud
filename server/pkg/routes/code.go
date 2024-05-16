package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patnaikankit/Code-Cloud/server/pkg/controllers"
)

func Code(ctx *gin.RouterGroup) {
	ctx.PUT("/code", controllers.CodeUpdate)
}
