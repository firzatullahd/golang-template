package main

import (
	"fmt"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/delivery/http/handler"
	"github.com/firzatullahd/cats-social-api/internal/delivery/http/route"
	"github.com/firzatullahd/cats-social-api/internal/repository"
	"github.com/firzatullahd/cats-social-api/internal/usecase"
)

func main() {
	conf := config.Load()
	fmt.Printf("%+v \n", conf)
	masterDB, replicaDB := config.InitializeDB(&conf.DB)

	repo := repository.NewRepository(masterDB, replicaDB)
	usecase := usecase.NewUsecase(conf, repo)
	handler := handler.NewHandler(usecase)
	route.Serve(conf, handler)
}
