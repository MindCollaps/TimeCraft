package router

import (
	"github.com/gin-gonic/gin"
	v1 "src/main/router/api/v1"
)

func InitRouter(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")

	v1.Handler(apiV1)

	return true
}
