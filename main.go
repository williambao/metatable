package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/williambao/metatable/config"
	"github.com/williambao/metatable/model"
	"github.com/williambao/metatable/router"
)

const Version = "0.1.0"

func main() {
	// 加载配置文件
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err.Error()))
	}

	// logrus.SetFormatter(&logrus.JSONFormatter{})
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)

		logrus.SetLevel(logrus.WarnLevel)

		logFile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err.Error())
		}
		logrus.SetOutput(logFile)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
	}
	// 数据库连接
	if err := model.OpenDB(cfg); err != nil {
		panic(fmt.Errorf("Fatal error connect to database : %s", err.Error()))
	}

	// 同步数据库结构
	if err := model.Migrate(); err != nil {
		panic(fmt.Errorf("Fatal error migrate tables to database : %s", err.Error()))
	}

	// if err := model.EnableCache(); err != nil {
	// 	panic(fmt.Errorf("Fatal error open database cache: %s", err.Error()))
	// }

	// // 缓存
	// _, err := cache.NewCache("memory", `{"interval":60}`)
	// if err != nil {
	// 	panic(fmt.Errorf("无法初始化缓存: %s", err.Error()))
	// }

	handler := router.Load(
		ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true),
		router.SetConfig(cfg),
	)
	listen := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := endless.NewServer(listen, handler)
	server.BeforeBegin = func(add string) {
		logrus.Debugf("Actual pid is %d", syscall.Getpid())
	}

	fmt.Printf("start app on: http://localhost:%d\n", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	log.Printf("Server on %s stopped\n", listen)

	os.Exit(0)
}
