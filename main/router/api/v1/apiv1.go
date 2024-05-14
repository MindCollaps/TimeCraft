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
	stygrpHandler(semestergroup)

	studygroup := rg.Group("/stygrp")
	sgrpHandler(studygroup)

	timeslot := rg.Group("/tsl")
	tslHandler(timeslot)

	roomconfig := rg.Group("/rmc")
	rmcHandler(roomconfig)

	roomspecialisation := rg.Group("rms")
	rmsHandler(roomspecialisation)
}
