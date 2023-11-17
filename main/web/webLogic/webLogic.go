package webLogic

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"src/main/database/models"
	"src/main/middleware"
)

func GetLogicData(c *gin.Context, path string) interface{} {
	//check if path exists in templateMap

	var data interface{} = nil
	if _, ok := templateMap[path]; ok {
		data = templateMap[path](c)
	} else {
		data = templateMap[""](c)
	}
	dat := make(map[string]interface{})

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	// Iterate over the fields of the struct
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		value := dataValue.Field(i).Interface()
		dat[field.Name] = value
	}

	c.Set("ignoreAuth", true)
	loggedIn := isLoggedIn(c)
	c.Set("ignoreAuth", false)
	dat["LoggedIn"] = loggedIn
	if loggedIn {
		data, exists := c.Get("user")
		var userData models.User
		if exists {
			userData = data.(models.User)
			dat["User"] = userData
		} else {
			dat["User"] = nil
		}

		data, exists = c.Get("userId")
		if exists {
			dat["UserId"] = data
		} else {
			dat["UserId"] = nil
		}
	} else {
	}

	return dat
}

type DefaultStruct struct {
}

func isLoggedIn(c *gin.Context) bool {
	middleware.LoginToken()(c)

	val, exists := c.Get("loggedIn")

	var isLoggedIn bool
	if exists {
		isLoggedIn = val.(bool)
	} else {
		isLoggedIn = false
	}

	return isLoggedIn
}

func index(c *gin.Context) any {
	return DefaultStruct{}
}

func defaultStruct(c *gin.Context) any {
	return DefaultStruct{}
}
