package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Printf("Server stop err: %v", err)
	} else {
		log.Printf("Server exit")
	}

	fmt.Println("main exited")

	// 关闭服务
	// lsof -i tcp:8080
	// kill pid
}

func Run(ctx context.Context) error {
	// 注册路由
	router := registerRoutes()

	// ...
	g, c := errgroup.WithContext(ctx)

	g.Go(func() error {
		return Server(c, "0.0.0.0:8080", router)
	})

	g.Go(func() error {
		return Server(c, "0.0.0.0:8081", http.DefaultServeMux)
	})

	g.Go(func() error {
		return Signal(c)
	})

	// inject error
	//g.Go(func() error {
	//	fmt.Println("inject")
	//	time.Sleep(5 * time.Second)
	//	fmt.Println("inject finish")
	//	return errors.New("inject error")
	//})

	return g.Wait()
}

// 启动服务
func Server(ctx context.Context, addr string, handler http.Handler) error {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Println("http")
	// 如果 s.Shutdown 执行成功之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-ctx.Done() // 等待 stop 信号
		fmt.Println("http ctx done")
		log.Printf("服务退出: %s\n", addr)
		s.Shutdown(context.TODO()) // 如果是 s.Shutdown(ctx), 有陷阱,shutdown没执行完, main goroutine就退出了.
	}()

	return s.ListenAndServe()
}

// 监听系统退出信号
func Signal(ctx context.Context) error {
	exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, len(exitSignals))
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, exitSignals...)
	for {
		fmt.Println("signal")
		select {
		case <-ctx.Done():
			fmt.Println("signal ctx done")
			return ctx.Err()
		case q := <-quit:
			fmt.Println("quit")
			return fmt.Errorf("收到的退出信号是: %v", q.String())
		}
	}

}

// registerRoutes ...
func registerRoutes() http.Handler {
	// 注册路由
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": http.StatusOK,
			"msg":  "hello!",
		})
		return
	})

	return router
}
