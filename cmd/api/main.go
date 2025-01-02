package main

import (
	"context"

	"github.com/firzatullahd/golang-template/config"
	"github.com/firzatullahd/golang-template/internal/user/adapter/mailersend"
	"github.com/firzatullahd/golang-template/internal/user/delivery/http/handler"
	"github.com/firzatullahd/golang-template/internal/user/delivery/http/route"
	"github.com/firzatullahd/golang-template/internal/user/repository"
	"github.com/firzatullahd/golang-template/internal/user/service"
	"github.com/firzatullahd/golang-template/utils/logger"
)

func main() {
	conf := config.Load()
	logger.NewLogger(conf.AppEnv)

	connDB, err := config.NewPsqlClient(conf)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	redisConn, err := config.NewRedisClient(context.Background(), conf)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	var userService service.Service
	{
		emailClient := mailersend.NewClient(conf.MailerSendAPIKey)
		repo := repository.NewRepository(connDB, connDB)
		userService = service.NewService(conf, repo, redisConn, emailClient)
	}

	handler := handler.NewHandler(&userService)
	route.Serve(conf, handler)
}
