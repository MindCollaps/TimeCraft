package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"src/main/database/models"
	"src/main/middleware"
)

const timeFormat = "2006-01-02 15:04:05"

type DashboardStruct struct {
	TimeTables string
}

func Dashboard(c *gin.Context) any {
	middleware.LoginToken()(c)

	if c.IsAborted() {
		return DashboardStruct{
			TimeTables: "[]",
		}
	}

	usr, ok := c.Get("user")
	if !ok {
		return DashboardStruct{
			TimeTables: "[]",
		}
	}

	user := usr.(models.User)

	timeTables, err := models.LoadTimeTables(c, user.StaredTimeTableIds)

	if err != nil {
		return DashboardStruct{
			TimeTables: "[]",
		}
	}

	timeTablesJson, err := json.Marshal(timeTables)

	return DashboardStruct{
		TimeTables: string(timeTablesJson),
	}
}
