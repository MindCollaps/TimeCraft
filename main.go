package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"src/main/core"
	"src/main/core/logChopper"
	"src/main/crypt"
	"src/main/database"
	"src/main/env"
	"src/main/router"
	"src/main/web/mail"
)

//go:embed main/web/*
var Files embed.FS

func main() {
	//File System setup
	_, err := Files.ReadDir("main/web")
	if err != nil {
		log.Println("Failed to read public files - this is likely a problem during compilation. Exiting...")
		return
	}

	env.Files = Files

	// command line arguments
	flag.BoolVar(&env.UNIX, "unix", false, "Run the server in unix mode")
	flag.BoolVar(&env.DEBUG, "debug", false, "Run the server in debug mode")
	flag.Parse()

	//Logger
	logChopper.LogChop()

	//Banner
	log.Println("\n" + env.BANNER + "\nTimeCraft" + "\nVersion: " + env.VERSION)

	//Environment setup
	envLocation := ".env"
	if env.UNIX {
		envLocation = "/etc/aso/.env"
	}

	if err := godotenv.Load(envLocation); err != nil {
		log.Println("No .env file found")
		return
	}

	//Crypt setup
	err = crypt.KeySetup()
	if err != nil {
		log.Println("Failed to setup keys")
		return
	}

	//Database setup
	database.InitDatabase()

	//Mail
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

	//Gin setup
	r := gin.Default()
	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	router.InitRouter(r)

	address := os.Getenv("ADDRESS")

	if address == "" {
		log.Println("No address set in .env file")
		address = "127.0.0.1:8080"
		log.Println("Defaulting to " + address)
	} else {
		log.Println("Listening to " + address)
	}

	err = r.Run(address)

	if err != nil {
		log.Fatal(err)
		return
	}
}
