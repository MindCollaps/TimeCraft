package v1

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"src/main/database"
	"src/main/database/models"
	"src/main/middleware"
	"strings"
)

// /api/v1/dev/...
func devHandler(cg *gin.RouterGroup) {
	cg.POST("/opt/model", middleware.TestingPurpose(), func(c *gin.Context) {
		type requestBody struct {
			Name string `json:"name" binding:"required"`
		}

		var req requestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Bad request"})
			log.Println(err)
			return
		}

		//get timetable by name
		var timetable models.TimeTable
		res := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"name": req.Name})
		if res.Err() != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "Timetable not found"})
			log.Println(res.Err())
			return
		}
		res.Decode(&timetable)

		tts, err := models.TimeTableToStruct(c, timetable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		var resData response
		resData.Name = tts.Name
		for _, day := range tts.Days {
			for _, slot := range day.TimeSlotIds {
				if slot.IsExam {
					resData.SpecialDays = append(resData.SpecialDays, specialDay{SpecialDayName: slot.Name, Day: day.Date, IsExam: true})
				} else if slot.IsHoliday {
					resData.SpecialDays = append(resData.SpecialDays, specialDay{SpecialDayName: slot.Name, Day: day.Date, IsHoliday: true})
				} else if slot.IsEvent {
					resData.SpecialDays = append(resData.SpecialDays, specialDay{SpecialDayName: slot.Name, Day: day.Date, IsEvent: true})
				} else if slot.IsReExamination {
					resData.SpecialDays = append(resData.SpecialDays, specialDay{SpecialDayName: slot.Name, Day: day.Date, IsReExamination: true})
				} else {
					name := slot.LecturerId.FirstName + " " + slot.LecturerId.SureName
					//strip space at end and start
					name = strings.TrimSpace(name)
					addToResponse(&resData, name, slot.LectureId.Name)
				}
			}
		}

		c.JSON(http.StatusOK, resData)
	})

	cg.GET("/opt/model", middleware.TestingPurpose(), func(c *gin.Context) {
		//get all timetable names
		col, err := database.MongoDB.Collection("TimeTable").Find(c, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		var timetables []string
		for col.Next(c) {
			var timetable models.TimeTable
			col.Decode(&timetable)
			timetables = append(timetables, timetable.Name)
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Fetched timetables", "timetables": timetables})
	})
}

type lecture struct {
	LectureName string `json:"lectureName"`
	Amount      int    `json:"amount"`
}

type lecturer struct {
	LecturerName string    `json:"lecturerName"`
	Lectures     []lecture `json:"lectures"`
}

type specialDay struct {
	SpecialDayName  string `json:"SpecialDayName"`
	Day             string `json:"day"`
	IsExam          bool   `json:"isExam"`
	IsHoliday       bool   `json:"isHoliday"`
	IsEvent         bool   `json:"isEvent"`
	IsReExamination bool   `json:"isReExamination"`
}

type response struct {
	Name        string       `json:"name"`
	Lecturers   []lecturer   `json:"lecturers"`
	SpecialDays []specialDay `json:"specialDay"`
}

func addToResponse(resData *response, lecturerName string, lectureName string) {
	for i, lecturer := range resData.Lecturers {
		if lecturer.LecturerName == lecturerName {
			for j, lecture := range lecturer.Lectures {
				if lecture.LectureName == lectureName {
					resData.Lecturers[i].Lectures[j].Amount++
					return
				}
			}
			resData.Lecturers[i].Lectures = append(resData.Lecturers[i].Lectures, lecture{LectureName: lectureName, Amount: 1})
			return
		}
	}
	resData.Lecturers = append(resData.Lecturers, lecturer{LecturerName: lecturerName, Lectures: []lecture{{LectureName: lectureName, Amount: 1}}})
}
