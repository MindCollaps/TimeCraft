package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"src/main/database"
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
	var LastChanged = time.Time{}
	for _, element := range data {

		if element.StudySubject != "" {
			fmt.Println("this is the header")
			LastChanged = getLastChanged(element.LastChanged)
			fmt.Println(LastChanged)
		} else {
			fmt.Println("this is the content")

			// make it iterable
			DaysToParse := element.Days
			v := reflect.ValueOf(DaysToParse)
			typeOfDays := v.Type()

			// iterate over the days
			for i := 0; i < v.NumField(); i++ {
				// get the day
				day := v.Field(i)
				// get the day name
				dayName := typeOfDays.Field(i).Name
				fmt.Println(dayName)

				// iterate over the lessons
				for _, lesson := range day.Interface().(Day).Lessons {
					startTime, endTime := getStartAndEndTime(lesson.Time)
					timeslot := models.TimeSlot{
						ID:              primitive.NewObjectID(),
						Name:            lesson.Name,
						LecturerId:      getLecturer(lesson.Lecturer),
						LectureId:       getLecture(element.StudySubject),
						TimeStart:       startTime,
						TimeEnd:         endTime,
						IsOnline:        lesson.IsOnline,
						IsReExamination: lesson.IsReExamination,
						IsExam:          lesson.IsExam,
						IsCancelled:     lesson.WasCanceled,
						WasMoved:        lesson.WasMoved,
						IsEvent:         lesson.IsEvent,
						RoomConfigId:    primitive.NilObjectID, // TODO: support rooms
						LastUpdated:     primitive.NewDateTimeFromTime(LastChanged),
					}
					saveTimeSlot(timeslot)
				}

			}

			// make a new struct from the database models for each day
			// save the struct to the database

		}
	}

	c.JSON(201, gin.H{"msg": "created"})
}

func getLastChanged(input string) time.Time {
	if strings.HasPrefix(input, "Stand: ") {
		input = strings.TrimPrefix(input, "Stand: ")
	}
	LastChanged, err := time.Parse("02.01.2006", input)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}
	} else {
		return LastChanged
	}
}

func getStartAndEndTime(lessonTime string) (primitive.DateTime, primitive.DateTime) {
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

func getLecturer(lecturer string) primitive.ObjectID {
	if lecturer == "" {
		return primitive.NilObjectID
	} else {
		var lecturerObj models.Lecturer

		err := database.MongoDB.Collection("Lecturer").FindOne(context.Background(), bson.M{
			"sureName": lecturer,
		}).Decode(&lecturerObj)

		if err != nil {
			fmt.Println("Error finding lecturer:", err)
			fmt.Println("Creating new lecturer")

			return saveLecturer(lecturer)
		} else {
			return lecturerObj.ID
		}
	}
}

func saveLecturer(lecturer string) primitive.ObjectID {
	lecturerObj := models.Lecturer{
		ID:       primitive.NewObjectID(),
		SureName: lecturer,
	}

	_, err := database.MongoDB.Collection("Lecturer").InsertOne(context.Background(), lecturerObj)
	if err != nil {
		fmt.Println("Error creating new lecturer:", err)
		return primitive.NilObjectID
	} else {
		return lecturerObj.ID
	}
}

func getLecture(lecture string) primitive.ObjectID {
	if lecture == "" {
		return primitive.NilObjectID
	} else {
		var lectureObj models.Lecture

		err := database.MongoDB.Collection("Lecture").FindOne(context.Background(), bson.M{
			"name": lecture,
		}).Decode(&lectureObj)

		if err != nil {
			fmt.Println("Error finding lecture:", err)
			fmt.Println("Creating new lecture")

			return saveLecture(lecture)
		} else {
			return lectureObj.ID
		}
	}
}

func saveLecture(lecture string) primitive.ObjectID {
	lectureObj := models.Lecture{
		ID:   primitive.NewObjectID(),
		Name: lecture,
	}

	_, err := database.MongoDB.Collection("Lecture").InsertOne(context.Background(), lectureObj)
	if err != nil {
		fmt.Println("Error creating new lecture:", err)
		return primitive.NilObjectID
	} else {
		return lectureObj.ID
	}
}

func saveTimeSlot(timeSlot models.TimeSlot) primitive.ObjectID {
	_, err := database.MongoDB.Collection("TimeSlot").InsertOne(context.Background(), timeSlot)
	if err != nil {
		fmt.Println("Error creating new timeSlot:", err)
		return primitive.NilObjectID
	} else {
		return timeSlot.ID
	}
}
