package v1

import "github.com/gin-gonic/gin"

func Handler(rg *gin.RouterGroup) {
	user := rg.Group("/usr")
	userHandler(user)

	prototype := rg.Group("/prt")
	prtHandler(prototype)

	timetable := rg.Group("/tbl")
	tblHandler(timetable)

	semestergroup := rg.Group("/sgrp")
	sgrpHandler(semestergroup)

	studentgroup := rg.Group("/stgrp")
	stgrpHandler(studentgroup)

	lecturegroup := rg.Group("lgrp")
	lgrpHandler(lecturegroup)

	devhandler := rg.Group("/dev")
	devHandler(devhandler)
}
