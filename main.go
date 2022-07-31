package main

import (
	"context"
	"github.com/augustazz/camellia/config"
	"github.com/augustazz/camellia/logger"
	"github.com/augustazz/camellia/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	//load config
	conf := config.GetSrvConfig()

	//setup logger
	logger.SetupLogger(ctx, conf.App.Name, conf.Log)

	s := server.Server{
		Port:        conf.Web.Port,
		Ctx:         ctx,
		AuthTimeout: time.Duration(conf.Conn.AuthTimeout) * time.Second, //s
		IdleTimeout: time.Duration(conf.Conn.IdleTimeout) * time.Minute, //min
	}
	go s.Start()

	t := make(chan os.Signal, 1)
	signal.Notify(t, os.Interrupt, syscall.SIGTSTP, syscall.SIGTERM, syscall.SIGINT)
	<-t

	//服务停止
	s.Close()
	logger.Info("server stopped...")
}
