package routes

import "github.com/gin-gonic/gin"

func Server(ctx *gin.Engine) {
	api := ctx.Group("/api")

	Code(api)
	Git(api)
}
