package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"src/main/env"
)

func TestingPurpose() gin.HandlerFunc {
	return func(c *gin.Context) {
		if env.TESTING {
			log.Println("!!! A testing function has been used because TESTING is enabled either in the .env file or in the command line arguments. This is insecure and should be disabled in production!")
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		}
	}
}
