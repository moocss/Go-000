## 学习笔记


## 作业

1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## 项目目录结构
```
.
├── api
│   └── v1
│       └── user
│           └── user.go
├── cmd
│   ├── myapp
│   │   ├── main.go
│   │   └── wire.go
│   └── server
│       ├── grpc
│       │   └── server.go
│       └── http
│           └── server.go
├── config
│   └── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── dto
│   │   └── user.go
│   ├── model
│   │   └── user.go
│   ├── repository
│   │   ├── mysql.go
│   │   └── user.go
│   └── service
│       └── user.go
└── pkg
    ├── conf
    │   └── conf.go
    └── errcode
        └── code.go

```