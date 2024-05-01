package main

import (
	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/repository"
	"github.com/firzatullahd/cats-social-api/internal/usecase"
)

func main() {
	conf := config.Load()
	masterDB, replicaDB := config.InitializeDB(&conf.DB)

	repo := repository.NewRepository(masterDB, replicaDB)
	usecase := usecase.NewUsecase(conf, repo)
}
