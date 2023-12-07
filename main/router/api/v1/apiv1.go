package v1

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
	user := rg.Group("/usr")

	userHandler(user)
}

func Handler2(rg *gin.RouterGroup) {
	timetable := rg.Group("/tbl")

	tblHandler(timetable)
}
