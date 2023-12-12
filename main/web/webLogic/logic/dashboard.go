package logic

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database/models"
)

type DashboardStruct struct {
	TimeTables []models.TimeTableStruct
}

func Dashboard(c *gin.Context) DashboardStruct {
	usr, ok := c.Get("user")
	if !ok {
		return DashboardStruct{
			TimeTables: []models.TimeTableStruct{},
		}
	}

	user := usr.(models.User)

	timeTables, err := LoadTimeTables(c, user.StaredTimeTableIds)

	if err != nil {
		return DashboardStruct{
			TimeTables: []models.TimeTableStruct{},
		}
	}

	return DashboardStruct{
		TimeTables: timeTables,
	}
}

func LoadTimeTables(c *gin.Context, timeTables []primitive.ObjectID) ([]models.TimeTableStruct, error) {
	var timeTablesStruct []models.TimeTableStruct
	for i := 0; i < len(timeTables); i++ {

	}
}
