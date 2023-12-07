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

func Handler3(rg *gin.RouterGroup) {
	semestergroup := rg.Group("/sgrp")

	sgrpHandler(semestergroup)
}

func Handler4(rg *gin.RouterGroup) {
	studygroup := rg.Group("/stygrp")

	stygrpHandler(studygroup)
}
