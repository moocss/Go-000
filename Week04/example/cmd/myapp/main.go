package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/moocss/example/api/v1/user"
	"github.com/moocss/example/cmd/server/grpc"
	xhttp "github.com/moocss/example/cmd/server/http"
	"golang.org/x/sync/errgroup"
)

func main() {
	/*
	// 初始化配置文件
	config := conf.NewConfig("./../../config/config.yaml")
	if config == nil {
		log.Fatal("配置文件读取失败")
	}

	// 连接数据库
	db := repository.NewDB(config)
	defer db.DB.Close()
	if db == nil {
		// 数据库无法连接直接panic
		log.Panicf("数据库连接失败")
	}

	log.Println("数据库连接成功")

	// 初始化 repositorys 和 services
	userRepo := repository.NewUserRepository(db)
	userSrv := service.NewUserService(userRepo)
    */

	srv := InitializeApp("./../../config/config.yaml")

	// 注册路由
	r := http.NewServeMux()
	user.NewUserHandler(r, srv)

	// run server
	if err := Run(context.Background(), r); err != nil {
		log.Printf("Server stop err: %v", err)
	} else {
		log.Printf("Server exit")
	}

}

func Run(ctx context.Context, mux *http.ServeMux) error {
	g, c := errgroup.WithContext(ctx)

	// http server
	g.Go(func() error {
		return HttpServer(c, mux)
	})

	// grpc server
	g.Go(func() error {
		return GrpcServer(c)
	})

	// signal
	g.Go(func() error {
		return Signal(c)
	})

	return g.Wait()
}

// HttpServer .
func HttpServer(ctx context.Context, mux *http.ServeMux) error {
	s := xhttp.NewServer(mux)
	fmt.Println("http")
	go func() {
		<-ctx.Done() // 等待 stop 信号
		fmt.Println("http ctx done")
		s.Stop(context.TODO())
	}()

	return s.Start()
}

// GrpcServer .
func GrpcServer(ctx context.Context) error {
	s := grpc.NewServer()
	fmt.Println("grpc")
	go func() {
		<-ctx.Done() // 等待 stop 信号
		fmt.Println("grpc ctx done")
		s.Stop(context.TODO())
	}()

	return s.Start()
}

// Signal .
func Signal(ctx context.Context) error {
	exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
	quit := make(chan os.Signal, len(exitSignals))
	signal.Notify(quit, exitSignals...)
	for {
		fmt.Println("signal")
		select {
		case <-ctx.Done():
			fmt.Println("signal ctx done")
			return ctx.Err()
		case q := <-quit:
			fmt.Println("quit")
			return fmt.Errorf("退出信号是: %v", q.String())
		}
	}
}
