package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"src/main/core"
	"src/main/crypt"
	"src/main/database"
	"src/main/mail"
	"src/main/router"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Println("No ..env file found")
	}

	err := crypt.KeySetup()
	if err != nil {
		log.Println("Failed to setup keys")
		return
	}

	database.InitDatabase()
	s := mail.InitMailer()
	if !s {
		disabled := os.Getenv("MAIL_DISABLED")
		if disabled == "true" {
			log.Println("Mail is disabled - ignoring error")
		} else {
			log.Fatal("Failed to setup mailer")
			return
		}
	}

	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	router.InitRouter(r)

	//set address
	address := os.Getenv("PORT")

	if address == "" {
		log.Println("No address set in .env file")
		address = ":8080"
		log.Println("Defaulting to " + address)
	} else {
		log.Println("Listening to " + address)
	}

	r.Run(address)
}
