package v1

import (
	"github.com/gin-gonic/gin"
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
type Lessons struct {
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
	Date    string    `json:"date"`
	Lessons []Lessons `json:"lessons"`
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
	c.JSON(201, gin.H{"msg": "created"})
}
