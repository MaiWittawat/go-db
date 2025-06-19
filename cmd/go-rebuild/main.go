package main

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	app_setup "go-rebuild/cmd/go-rebuild/setup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// ------------------------------ Setup Config ------------------------------
	appcore_config.InitConfigurations()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	gin.SetMode(gin.ReleaseMode)

	if appcore_config.Config.Mode == "develop" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}



	// ------------------------------ Init db ------------------------------
	appClients, err := app_setup.InitAppClients()
	if err != nil {
		log.Panicf("Failed to initialize 3rd party clients: %v", err)
	}



	// ------------------------------ Start service ------------------------------
	appServices := app_setup.BuildApplicationServices(appClients)
	router := gin.Default()



	// ------------------------------ Register API ------------------------------
	app_setup.APIRoutes(router, appServices)



	// ------------------------------ Start consume -----------------------------
	app_setup.StartConsumers(appServices.MQBroker)



	// ------------------------------ Start server ------------------------------
	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("[server]: Server closed gracefully")
			} else {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
	log.Info("[server]: server start at port:3000")



	// ------------------------------ Shutdown ------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("[Signal]: shutdown signal received")
	app_setup.GracefulShutdown(shutdownCtx, server, appClients)
}
