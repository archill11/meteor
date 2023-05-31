package main

import (
	"log"
	"meteor/internal/api/controllers"
	"meteor/internal/api/server"
	"meteor/internal/config"
	"meteor/internal/provider"
	"meteor/internal/service"
	"os"
	"os/signal"
	"syscall"

	jsoniter "github.com/json-iterator/go"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

//	@title			meteor
//	@version		0.1
//	@description	description
//	@host			meteor-api.onrender.com
//	@BasePath

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("can't load config:", err)
	}

	// zapCfg := zap.NewProductionConfig()
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.OutputPaths = []string{"stderr"}
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal("can't init logger", err)
	}
	defer logger.Sync()

	json := jsoniter.ConfigFastest

	provider := provider.New(config.Provider, logger, json)

	service := service.New(config.Service, provider, logger, json)

	b := controllers.NewBaseController(service)
	r := server.NewRouters(struct {
		*controllers.MainController
		*controllers.MarketplaceController
	}{
		controllers.NewMainControllerFromBase(b),
		controllers.NewMarketplaceControllerFromBase(b),
	})

	HTTPServer := server.NewServer(config.Server, r)

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("[StartServer]", zap.String("Address", config.Server.ListenAddress))
		serverErrors <- HTTPServer.ListenAndServe(config.Server.ListenAddress)
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Fatal("Server error", zap.Error(err), zap.String("action", "FATAL_SERVICE"))

	case <-osSignals:
		logger.Info("Shutdown...")
		if err := HTTPServer.Shutdown(); err != nil {
			logger.Info("Gracefull shutdown in  not completed", zap.Error(err), zap.String("action", "EXIT_SERVICE"))
		}
		os.Exit(0)
	}

}
