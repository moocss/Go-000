//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/moocss/example/internal/repository"
	"github.com/moocss/example/internal/service"
	"github.com/moocss/example/pkg/conf"
)

func InitializeApp(path string) service.UserService {
	wire.Build(conf.NewConfig, repository.NewDB, repository.NewUserRepository, service.NewUserService)
	return nil
}
