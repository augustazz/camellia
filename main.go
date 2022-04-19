package main

import (
	"camellia/config"
	"camellia/logger"
	"camellia/server"
	"os"
)

var defaultConfigFilePath = "resources/config.yml"

func main() {
	profiles := os.Getenv("profiles")
	configPath := os.Getenv("configPath")
	if configPath == "" {
		configPath = defaultConfigFilePath
	}

	//load config
	conf := config.ServerConfig{}
	conf.LoadConfig(profiles, configPath)

	//setup logger
	logger.SetupLogger(conf.App.Name, conf.Log)

	s := server.Server{
		Port: conf.Web.Port,
	}
	s.Start()
}
