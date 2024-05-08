package webLogic

import (
	"github.com/gin-gonic/gin"
	"src/main/web/webLogic/logic"
)

var templateMap = map[string]func(c *gin.Context) any{
	".":         index,
	"dashboard": logic.Dashboard,
	"":          defaultStruct,
	// Add more entries as needed
}
