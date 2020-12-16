// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/moocss/example/internal/repository"
	"github.com/moocss/example/internal/service"
	"github.com/moocss/example/pkg/conf"
)

// Injectors from wire.go:

func InitializeApp(path string) service.UserService {
	confConf := conf.NewConfig(path)
	db := repository.NewDB(confConf)
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	return userService
}
