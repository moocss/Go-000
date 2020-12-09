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

	// 关闭服务
	// lsof -i tcp:8080
	// kill pid
}

func Run(ctx context.Context) error {
	// 注册路由
	router := registerRoutes()

	// ...
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return Server(ctx, "0.0.0.0:8080", router)
	})

	g.Go(func() error {
		return Server(ctx, "0.0.0.0:8081", http.DefaultServeMux)
	})

	g.Go(func() error {
		return Signal(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	return g.Wait()
}

// 启动服务, 优化过的
func Server(ctx context.Context, addr string, handler http.Handler) error {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// 如果 s.Shutdown 执行成功之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-ctx.Done()
		log.Printf("服务退出: %s\n", addr)
		s.Shutdown(ctx)
	}()

	return s.ListenAndServe()
}

/*
// 启动服务
func Server(ctx context.Context, addr string, handler http.Handler) error {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	done := make(chan error)

	go func() {
		defer close(done)
		log.Printf("服务启动: %s\n", addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务启动失败: %v", err)
			// 发送服务器错误
			done <- err
		}
	}()

	select {
	case <-ctx.Done(): // 正常关机
	err := s.Shutdown(ctx)
		fmt.Printf("Server shutdown with error: %v\n", err)
		return s.Shutdown(ctx)
	case err := <-done: // 接收服务器错误
		return err
	}
}
*/

// 监听系统退出信号
func Signal(ctx context.Context) error {
	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-ctx.Done():
		signal.Reset()
		return nil
	case q := <-quit:
		return fmt.Errorf("收到的退出信号是: %v", q.String())
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
