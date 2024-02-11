package main

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"attendit/backend/routes"
	"attendit/backend/services"
)

// @title GoLang Rest API Starter Doc
// @version 1.0
// @description GoLang - Gin - RESTful - MongoDB - Redis
// @termsOfService https://swagger.io/terms/

// @contact.name Ebubekir Yiğit
// @contact.url https://github.com/ebubekiryigit
// @contact.email ebubekiryigit6@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer-Token
func main() {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
		return
	}

	time.Local = location

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Println(err)
		return
	}

	services.LoadConfig()
	services.InitMongoDB()
	services.CheckRedisConnection()

	j, err := s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(23, 59, 59),
			),
		),
		gocron.NewTask(
			services.CheckOutAllAttendances,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}

	routes.InitGin()
	router := routes.New()

	server := &http.Server{
		Addr:         services.Config.ServerAddr + ":" + services.Config.ServerPort,
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			return
		}
	}()

	s.Start()

	nextTime, err := j.NextRun()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Next Cron Job run time: %v\n", nextTime)

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 15 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.Shutdown(); err != nil {
		log.Println("Scheduler Shutdown:", err)
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
