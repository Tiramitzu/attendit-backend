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
	services.LoadConfig()
	services.InitMongoDB()
	services.CheckRedisConnection()

	routes.InitGin()
	router := routes.New()

	server := &http.Server{
		Addr:         services.Config.ServerAddr + ":" + services.Config.ServerPort,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      router,
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.NewJob(
		gocron.WeeklyJob(
			1,
			gocron.NewWeekdays(time.Saturday),
			gocron.NewAtTimes(
				gocron.NewAtTime(1, 30, 0),
				gocron.NewAtTime(12, 0, 30),
			),
		),
		gocron.NewTask(
			func() {
				_, err := services.DeleteUnVerifiedUsers()
				if err != nil {
					log.Println(err)
					return
				}

				log.Println("Unverified Users Deleted.")
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Cron Job Started.")
	s.Start()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 15 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down cron job...")

	if err := s.Shutdown(); err != nil {
		log.Fatal("Cron Job Shutdown:", err)
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
