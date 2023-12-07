package v1

import (
	"github.com/gin-gonic/gin"
)

func tblHandler(cg *gin.RouterGroup) {
	//    /api/v1/tbl/...
	cg.POST("/create", func(c *gin.Context) {
		//check body for username and password
		var requestBody struct {
			Id   string `json:"id" binding:"required"`
			Name string `json:"name" binding:"required"`
			Days string `json:"days" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		id := requestBody.Id
		name := requestBody.Name
		days := requestBody.Days

		var existingTbl models.T
	})

}
