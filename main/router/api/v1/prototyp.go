package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"src/main/database/models"
	"strings"
	"time"
)

type ExcelJson []struct {
	StudySubject  string `json:"study_subject"`
	SemesterGroup string `json:"semester_group"`
	Semester      string `json:"semester"`
	SemesterYear  string `json:"semester_year"`
	LastChanged   string `json:"last_changed"`
	Days          Days   `json:"days"`
	Calendarweek  int    `json:"calendarweek"`
	StartRow      int    `json:"start_row"`
	EndRow        int    `json:"end_row"`
}
type Lesson struct {
	Time            string `json:"time"`
	Name            string `json:"name"`
	IsOnline        bool   `json:"isOnline"`
	IsReExamination bool   `json:"isReExamination"`
	IsExam          bool   `json:"isExam"`
	WasCanceled     bool   `json:"wasCanceled"`
	WasMoved        bool   `json:"wasMoved"`
	Lecturer        string `json:"lecturer"`
	IsEvent         bool   `json:"isEvent"`
}
type Day struct {
	Date    string   `json:"date"`
	Lessons []Lesson `json:"lessons"`
}

type Days struct {
	Montag     Day `json:"Montag"`
	Dienstag   Day `json:"Dienstag"`
	Mittwoch   Day `json:"Mittwoch"`
	Donnerstag Day `json:"Donnerstag"`
	Freitag    Day `json:"Freitag"`
	Samstag    Day `json:"Samstag"`
}

func prtHandler(cg *gin.RouterGroup) {
	cg.POST("/import/excel", func(c *gin.Context) {
		// get file from body
		// save file to temp folder
		// excecute Python script and generate json file
		// parse json file and save to database
		c.JSON(501, gin.H{"msg": "not implemented yet"})
		return
	})

	cg.POST("/import/exceljson", func(c *gin.Context) {
		// get json data from body
		var requestBody struct {
			JsonData ExcelJson `json:"json_data"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		parse_json(requestBody.JsonData, c)
	})
}

func parse_json(data ExcelJson, c *gin.Context) {
	// parse json file and save to database
	for _, element := range data {

		if element.StudySubject != "" {
			fmt.Println("this is the header")
		} else {
			fmt.Println("this is the content")

			// make it iterable
			DaysToParse := element.Days
			v := reflect.ValueOf(DaysToParse)
			typeOfDays := v.Type()

			var TimeSlots []models.TimeSlot
			// iterate over the days
			for i := 0; i < v.NumField(); i++ {
				// get the day
				day := v.Field(i)
				// get the day name
				dayName := typeOfDays.Field(i).Name
				fmt.Println(dayName)

				// iterate over the lessons
				for _, lesson := range day.Interface().(Day).Lessons {
					// get the lesson

					// TODO:
					//	- get the id of the lectuerer
					// 		- lesson has no lectuerer --> return primitive.NilObjectID
					// 		- lecturer doesn't exist in DB yet --> create new lecturer and return the id
					// 		- lecturer exists in DB --> return the id
					//
					//	- same for the LectureID and RoomConfigID

					startTime, endTime := parse_time(lesson.Time)
					timeslot := models.TimeSlot{
						ID:              primitive.NewObjectID(),
						Name:            lesson.Name,
						LecturerID:      primitive.NilObjectID,
						LectureID:       primitive.ObjectID{},
						TimeStart:       startTime,
						TimeEnd:         endTime,
						IsOnline:        lesson.IsOnline,
						IsReExamination: lesson.IsReExamination,
						IsExam:          lesson.IsExam,
						IsCancelled:     lesson.WasCanceled,
						WasMoved:        lesson.WasMoved,
						IsEvent:         lesson.IsEvent,
						RoomConfigID:    primitive.ObjectID{},
					}
					TimeSlots = append(TimeSlots, timeslot)
				}

			}

			// make a new struct from the database models for each day
			// save the struct to the database

		}
	}

	c.JSON(201, gin.H{"msg": "created"})
}

func parse_time(lessonTime string) (primitive.DateTime, primitive.DateTime) {
	timeRange := strings.Split(lessonTime, "-")
	startTimeStr := timeRange[0]
	endTimeStr := timeRange[1]

	layout := "15:04"
	startTime, err := time.Parse(layout, startTimeStr)
	endTime, err := time.Parse(layout, endTimeStr)

	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	return primitive.NewDateTimeFromTime(startTime), primitive.NewDateTimeFromTime(endTime)

}
