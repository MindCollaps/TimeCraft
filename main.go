package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"src/main/core"
	"src/main/crypt"
	"src/main/database"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	err := crypt.KeySetup()
	if err != nil {
		log.Println("Failed to setup keys")
		return
	}

	database.InitDatabase()

	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	//set address
	address := os.Getenv("ADDRESS")

	if address == "" {
		log.Println("No address set in .env file")
		address = ":8080"
		log.Println("Defaulting to " + address)
	} else {
		log.Println("Listening on " + address)
	}

	r.Run(address)
}
