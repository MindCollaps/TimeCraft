package v1

import (
	"github.com/gin-gonic/gin"
)

func userHandler(cg *gin.RouterGroup) {
	//    /api/v1/usr/login
	cg.GET("/login", func(c *gin.Context) {

	})

	cg.POST("/register", func(c *gin.Context) {
		var requestBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		if err := c.ShouldBindJSON(requestBody); err != nil {
			c.JSON(400, gin.H{
				"msg":    "bad request",
				"status": 400,
			})

			return
		}

	})
}
