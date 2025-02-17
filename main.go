package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/test/library-app/docs"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/handler"
	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/store"
)

// @title 		Library App
// @version 	1.0
// @description Handles the digital library operation
// @host 		localhost:3000
// @basePath 	/api/v1
func main() {
	// loads config if any error in reading config panics the appl
	err := config.LoadConfig()
	if err != nil {
		log.Panicf("%v", err)
	}

	// configures logger for an app
	logger.InitLogger()
	logger.Infof("Hello this is library-app")

	// initializing the gin router
	router := gin.Default()

	// Initializing the store
	store, err := store.NewStore()
	if err != nil {
		logger.Panicf("failed to initialize store. Error:%v", err)
	}
	defer store.Close()

	// Actual handler to handles the requests
	handler := handler.NewHandler(store)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	bookRouter := router.Group("/api/v1")
	{
		bookRouter.GET("/hello", handler.Hello)
		bookRouter.GET("/book", handler.GetAllBooks)
		bookRouter.GET("/book/:title", handler.GetBook)
		bookRouter.POST("/borrow", handler.BorrowBook)
		bookRouter.POST("/extend/:id", handler.ExtendLoan)
		bookRouter.POST("/return/:id", handler.ReturnBook)
	}

	// Attaching the request handlers, port etc to the server
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", config.CommonConfig.ServicePort),
		Handler:      router,
		ReadTimeout:  time.Duration(config.CommonConfig.ReadTimeoutInSec) * time.Second,
		WriteTimeout: time.Duration(config.CommonConfig.WriteTimeoutInSec) * time.Second,
		IdleTimeout:  time.Duration(config.CommonConfig.IdleTimeoutInSec) * time.Second,
	}
	// holds the server related errors
	serverErros := make(chan error, 1)
	// starting a http server with seperate go routine
	go func() {
		logger.Infof("Server is up and running on: %v", config.CommonConfig.ServicePort)
		if err := server.ListenAndServe(); err != nil {
			serverErros <- err
		}
	}()

	// handling the graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// main go rotines hangs here till any of the errors for starting the server or shutdown signal
	select {
	case err := <-serverErros:
		logger.Errorf("Failed to start the server. Error: %v", err)
	case sig := <-stop:
		logger.Infof("Shutting down the app with signal: %v", sig)
	}
	logger.Infof("Server is shutting down")

	// proceeding further to cleanup any resources like servers, db connections etc
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Server handles on going requests, forcing it to shutdown after timeout to avoid hanging over some indefinite requests
	if err := server.Shutdown((ctx)); err != nil {
		logger.Errorf("Failed to shutdown the server properly. Error: %v", err)
	}
	logger.Infof("Server exited gracefully")
}
