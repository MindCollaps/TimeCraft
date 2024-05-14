package logChopper

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"src/main/env"
)

func LogChop() {
	filePath := "timecraft.log"
	if env.UNIX {
		filePath = "/var/logChopper/timecraft/timecraft.log"
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create a multi-writer that writes to both os.Stdout and the logChopper file
	mw := io.MultiWriter(os.Stdout, file)

	// Set logChopper output to the multi-writer
	log.SetOutput(mw)

	// Set Gin's debug mode and writer
	if env.DEBUG {
		log.Println("!!! Server is running in debug mode!")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if env.TESTING {
		log.Println("!!! Server is running in testing mode - please disable this in production because insecure functions are exposed!")
	}

	gin.DefaultWriter = mw
}
