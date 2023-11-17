package webLogic

import "github.com/gin-gonic/gin"

var templateMap = map[string]func(c *gin.Context) any{
	".": index,
	"":  defaultStruct,
	// Add more entries as needed
}
