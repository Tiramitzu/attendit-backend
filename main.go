package main

import (
	"context"
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

	log.Println("Cron Job Started.")

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
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
