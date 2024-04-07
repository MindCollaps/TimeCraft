package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Days          []Day  `json:"days"`
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
	IsHoliday       bool   `json:"isHoliday"`
}
type Day struct {
	Date    string   `json:"date"`
	Lessons []Lesson `json:"lessons"`
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

	cg.POST("/import/json", func(c *gin.Context) {
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
	var LastChanged primitive.DateTime
	var Name string
	var TimeTableDays []primitive.ObjectID
	var TimeTable models.TimeTable

	for _, element := range data {

		if element.StudySubject != "" {
			fmt.Println("parsing the header")
			LastChanged = getLastChanged(element.LastChanged)
			Name = fmt.Sprint(element.StudySubject, " ", element.SemesterGroup, " ", strings.Replace(element.SemesterYear, " ", "", 1))
		} else {
			fmt.Println(fmt.Sprintf("parsing calendarweek %d", element.Calendarweek))

			// make it iterable
			DaysToParse := element.Days
			for _, day := range DaysToParse {
				dayDate := day.Date

				var TimeTableDay models.TimeTableDay
				var TimeSlotIds []primitive.ObjectID
				TimeTableDay.ID = primitive.NewObjectID()
				TimeTableDay.Date = convertToDateTime("2006-01-02 00:00:00", dayDate)
				TimeTableDay.LastUpdated = LastChanged

				// iterate over the lessons
				for _, lesson := range day.Lessons {
					// skip empty time slots
					if strings.HasPrefix(lesson.Name, "no lesson") {
						continue
					}

					startTime, endTime := getStartAndEndTime(lesson.Time)
					timeslot := models.TimeSlot{
						ID:              primitive.NewObjectID(),
						Name:            lesson.Name,
						LecturerId:      getLecturer(lesson.Lecturer),
						LectureId:       getLecture(lesson),
						TimeStart:       startTime,
						TimeEnd:         endTime,
						IsOnline:        lesson.IsOnline,
						IsReExamination: lesson.IsReExamination,
						IsExam:          lesson.IsExam,
						IsCancelled:     lesson.WasCanceled,
						WasMoved:        lesson.WasMoved,
						IsEvent:         lesson.IsEvent,
						IsHoliday:       lesson.IsHoliday,
						RoomConfigId:    primitive.NilObjectID, // TODO: support rooms
						LastUpdated:     LastChanged,
					}

					TimeSlotIds = append(TimeSlotIds, saveTimeSlot(timeslot))
				}

				TimeTableDay.TimeSlotIds = TimeSlotIds
				id := saveTimeTableDay(TimeTableDay)
				TimeTableDays = append(TimeTableDays, id)
			}
		}
	}

	TimeTable.ID = primitive.NewObjectID()
	TimeTable.Name = Name
	TimeTable.Days = TimeTableDays

	id := saveTimeTable(TimeTable)
	fmt.Println(id)

	c.JSON(201, gin.H{"msg": "created"})
}

func getLastChanged(input string) primitive.DateTime {
	if strings.HasPrefix(input, "Stand: ") {
		input = strings.TrimPrefix(input, "Stand: ")
	}

	return convertToDateTime("02.01.2006", input)
}

func getStartAndEndTime(lessonTime string) (primitive.DateTime, primitive.DateTime) {
	timeRange := strings.Split(lessonTime, "-")
	startTimeStr := timeRange[0]
	endTimeStr := timeRange[1]

	startTimeStr = "2020-01-01 " + startTimeStr + ":00"
	endTimeStr = "2020-01-01 " + endTimeStr + ":00"

	startTime := convertToDateTime(time.DateTime, startTimeStr)
	endTime := convertToDateTime(time.DateTime, endTimeStr)

	return startTime, endTime

}

func convertToDateTime(layout string, input string) primitive.DateTime {
	//set timezone to local
	loc, _ := time.LoadLocation("Europe/Berlin")
	parsedTime, err := time.ParseInLocation(layout, input, loc)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	return primitive.DateTime(primitive.NewDateTimeFromTime(parsedTime))
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
			fmt.Println(fmt.Sprintf("Creating new lecturer '%s'", lecturer))
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

func getLecture(lecture Lesson) primitive.ObjectID {
	if lecture.Name == "" || lecture.IsEvent || lecture.IsExam || lecture.IsReExamination || lecture.IsHoliday || strings.HasPrefix(lecture.Name, "no lesson") {
		return primitive.NilObjectID
	} else {
		var lectureObj models.Lecture

		err := database.MongoDB.Collection("Lecture").FindOne(context.Background(), bson.M{
			"name": lecture.Name,
		}).Decode(&lectureObj)

		if err != nil {
			fmt.Println(fmt.Sprintf("Creating new lecture '%s'", lecture.Name))
			return saveLecture(lecture.Name)
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

func saveTimeTableDay(timeTableDay models.TimeTableDay) primitive.ObjectID {
	_, err := database.MongoDB.Collection("TimeTableDay").InsertOne(context.Background(), timeTableDay)
	if err != nil {
		fmt.Println("Error creating new timeTableDay:", err)
		return primitive.NilObjectID
	} else {
		return timeTableDay.ID
	}
}

func saveTimeTable(timeTable models.TimeTable) primitive.ObjectID {
	_, err := database.MongoDB.Collection("TimeTable").InsertOne(context.Background(), timeTable)
	if err != nil {
		fmt.Println("Error creating new timeTable:", err)
		return primitive.NilObjectID
	} else {
		return timeTable.ID
	}
}
