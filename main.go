package main

import (
	"camellia/config"
	"camellia/logger"
	"camellia/server"
	"context"
	"os"
	"os/signal"
	"syscall"
)

var defaultConfigFilePath = "resources/config.yml"

func main() {
	ctx := context.Background()

	profiles := os.Getenv("profiles")
	configPath := os.Getenv("configPath")
	if configPath == "" {
		configPath = defaultConfigFilePath
	}

	//load config
	conf := config.ServerConfig{}
	conf.LoadConfig(ctx, profiles, configPath)

	//setup logger
	logger.SetupLogger(ctx, conf.App.Name, conf.Log)

	s := server.Server{
		Port: conf.Web.Port,
		Ctx: ctx,
	}
	go s.Start()


	t := make(chan os.Signal, 1)
	signal.Notify(t, os.Interrupt, syscall.SIGTSTP, syscall.SIGTERM, syscall.SIGINT)
	<-t

	//服务停止
	s.Close()
	logger.Info("server stopped...")
}
