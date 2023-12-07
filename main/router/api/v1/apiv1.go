package v1

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
	user := rg.Group("/usr")
	timetable := rg.Group("/tbl")
	semestergroup := rg.Group("/sgrp")
	studygroup := rg.Group("/stygrp")
	userHandler(user)
	tblHandler(timetable)
	stygrpHandler(studygroup)
	sgrpHandler(semestergroup)

}
