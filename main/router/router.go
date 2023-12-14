package router

import (
	"github.com/gin-gonic/gin"
	v1 "src/main/router/api/v1"
)

func InitRouter(r *gin.Engine) bool {
	apiV1 := r.Group("/api/v1")

	v1.Handler(apiV1)

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("auth", "", 0, "/", "", false, false)
		c.Redirect(302, "/")
	})

	return true
}
