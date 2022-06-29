package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controllers"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"go.uber.org/zap"
)

// go web 开发通用的脚手架模板

// @title bluebell项目接口文档
// @version 1.0
// @description Go web开发进阶项目实战课程bluebell

// @contact.name liwenzhou
// @contact.url http://www.liwenzhou.com

// @host 127.0.0.1:8081
// @BasePath /api/v1
func main() {

	// 以命令行的方式指定配置文件
	//if len(os.Args) < 2 {
	//	fmt.Println("need config file....")
	//	return
	//}
	//fileName := os.Args[1]

	var fileName2 string
	flag.StringVar(&fileName2, "f", "config.yaml", "config file")
	flag.Parse()
	fileName := fileName2
	fmt.Printf("config file: %s\n", fileName)
	//	1.加载配置
	err := settings.Init(fileName)
	if err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//	2.初始化日志
	err = logger.Init(settings.Conf.LogConfig, settings.Conf.AppConfig.Mode)
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	zap.L().Debug("logger init success.... ")

	defer zap.L().Sync()

	//	3.初始化mysql连接
	err = mysql.Init(settings.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	//	4.初始化redis连接
	err = redis.Init(settings.Conf.RedisConfig)
	if err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	err = snowflake.Init(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineID)
	if err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}

	//	5.注册路由
	zap.L().Info(settings.Conf.Name)
	r := routes.SetUp(settings.Conf.AppConfig.Mode)

	//	6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}
	fmt.Printf("web server %s running 127.0.0.1:%d \n", settings.Conf.AppConfig.Name, settings.Conf.AppConfig.Port)
	go func() {
		err2 := srv.ListenAndServe()
		if err2 != nil && err2 != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err2)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		zap.L().Fatal("server shutdown ", zap.Error(err))
	}
	zap.L().Info("server exiting...")

	//r.Run()

}
