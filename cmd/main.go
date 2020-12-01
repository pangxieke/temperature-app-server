package main

import (
	"log"
	"net/http"
	"os"
	"temperature/config"
	"temperature/middlewares"
	"temperature/models"
	"temperature/routers"
	"time"
)

func main() {
	Env := os.Getenv("GO_ENV")
	if Env != ""{
		logPath := "var/log/temperature-" + time.Now().Format("200601") + ".log"
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Panic(err)
		}
		log.SetOutput(logFile)
	}


	if err := config.Init(); err != nil {
		log.Panic(err)
	}
	if err := models.Init(); err != nil {
		log.Panic(err)
	}
	log.Println("server start:8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.Handler(routers.Router())))
}
