package v1

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
	user := rg.Group("/usr")
	timetable := rg.Group("/tbl")
	userHandler(user)
	tblHandler(timetable)
}
