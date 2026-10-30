// Command thsdk-server runs a THSDK client as an HTTP service.
//
// The same code is importable as a library: github.com/angYuHao/go-thsdk for
// direct calls and github.com/angYuHao/go-thsdk/httpserver to mount the routes
// inside an existing application.
//
// Usage:
//
//	thsdk-server [flags]
//
// Endpoints:
//
//	GET  /healthz             liveness and auth state
//	GET  /v1/methods          the exposed method catalog
//	POST /v1/call/<method>    one call, body is the method's params object
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
	"github.com/angYuHao/go-thsdk/httpserver"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "thsdk-server:", err)
		os.Exit(1)
	}
}

func run() error {
	addr := flag.String("addr", "127.0.0.1:8765", "HTTP listen address")
	runtimePath := flag.String("runtime", "", "path to hq.dll (default: auto-discover)")
	sessionDir := flag.String("session-dir", "", "directory holding account.session (default: working directory)")
	allowWrites := flag.Bool("allow-writes", false, "expose the 10 account-mutation methods")
	noTemporary := flag.Bool("no-temporary", false, "never fall back to the bundled trial accounts")
	callTimeout := flag.Duration("call-timeout", 30*time.Second, "per-request upstream budget when the client sends no deadline")
	authTimeout := flag.Duration("auth-timeout", 60*time.Second, "login handshake budget")
	shutdownTimeout := flag.Duration("shutdown-timeout", 10*time.Second, "graceful shutdown budget")
	flag.Parse()

	if *allowWrites {
		log.Println("警告：写方法已启用，调用会修改账号的云端自选股")
	}

	client, err := thsdk.New(thsdk.Options{
		RuntimePath:             *runtimePath,
		SessionDir:              *sessionDir,
		DisableTemporaryAccount: *noTemporary,
		AuthTimeout:             *authTimeout,
	})
	if err != nil {
		return fmt.Errorf("加载运行组件失败: %w", err)
	}
	defer client.Close()

	authCtx, cancelAuth := context.WithTimeout(context.Background(), *authTimeout)
	err = client.Auth(authCtx, thsdk.Credentials{})
	cancelAuth()
	if err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}
	log.Printf("认证完成，开始监听 %s", *addr)

	handler, err := httpserver.New(httpserver.Config{
		Client:      client,
		AllowWrites: *allowWrites,
		CallTimeout: *callTimeout,
		Logger:      log.Printf,
	})
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           handler.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serveErr:
		return err
	case <-stop:
		log.Println("收到退出信号，正在关闭")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("关闭失败: %w", err)
	}
	return nil
}
