package logic

import (
	"github.com/gin-gonic/gin"
	"src/main/database/models"
)

const timeFormat = "2006-01-02 15:04:05"

type DashboardStruct struct {
	TimeTables []models.TimeTableStruct
}

func Dashboard(c *gin.Context) any {
	usr, ok := c.Get("user")
	if !ok {
		return DashboardStruct{
			TimeTables: []models.TimeTableStruct{},
		}
	}

	user := usr.(models.User)

	timeTables, err := models.LoadTimeTables(c, user.StaredTimeTableIds)

	if err != nil {
		return DashboardStruct{
			TimeTables: []models.TimeTableStruct{},
		}
	}

	return DashboardStruct{
		TimeTables: timeTables,
	}
}
