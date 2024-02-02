package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nullptr-z/gin-template/router"
	"github.com/nullptr-z/gin-template/settings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	if err := settings.InitViperConfig(); err != nil {
		fmt.Println("settings.Init", err)
		return
	}
	// 2. 初始化日志
	if err := settings.InitLogger(); err != nil {
		fmt.Println("settings.Init", err)
		return
	}
	defer zap.L().Sync()
	// 3. 初始化 Sql DB
	if err := settings.InitSqlDB(); err != nil {
		fmt.Println("settings.Init", err)
		return
	}
	defer settings.CloseDB()
	// 4. 初始化 Redis
	if err := settings.InitRedis(); err != nil {
		fmt.Println("settings.Init", err)
		return
	}
	defer settings.CloseRedis()
	// 5. 注册服务路由
	g := router.Setup()
	// 6. 启动服务（优雅停机）
	startService(g)
}

func startService(g *gin.Engine) {
	host := viper.GetString("app.host")
	port := viper.GetInt("app.port")
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{Addr: addr, Handler: g}
	go func() {
		fmt.Print("\n Listening on: http://", addr, "\n\n")
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Error("Listening:", zap.Error(err))
		}
	}()

	// 等待系统中断信号来关闭服务，为关闭服务设置一个 5 秒的超时
	// kill default syscaLL.SIGTERM
	// kill -2  syscaLL.SIGINT 我们常用的ctrl+C就是触发系统 SIGINT 信号
	// kill -9  syscalL.SIGKILL 不能被捕获，所以不需要添加它
	quit := make(chan os.Signal, 1)
	// signal.Notify 会把收到的 syscalL.SIGINT 或 syscaLL.SIGTERM 信号转发给 quit
	<-quit // 阻塞等待关闭信号
	zap.L().Info("Shutdown Server...")
	// 定时 5 的Chan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 延迟 5 秒，处理还未完成的请求扫尾，然后优雅停机
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Service exiting")
}
