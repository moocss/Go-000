package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moocss/example/api"
	"github.com/moocss/example/internal/repository"
	"github.com/moocss/example/internal/service"
)

func main()  {
	port := 4000

	// 连接数据库
	db, err := service.OpenConnection()
	if err != nil {
		// 数据库无法连接直接panic
		log.Panicf("数据库连接失败: %s", err.Error())
	}

	log.Println("数据库连接成功")

	defer db.Close()

	// 初始化 repositorys 和 services
	userRepo := repository.NewUserRepository(db)
	userSrv := service.NewUserService(userRepo)

	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
		// 省略, 性能优化配置
	}

	// 路由
	api.NewHandler(userSrv)

	log.Printf("服务启动, 端口为: %d\n", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	defer Shutdown(server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-quit))
	log.Println("服务停止")
}

func Shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("无法关闭服务:", err)
	}
	log.Println("服务退出")
}
