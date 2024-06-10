package main

import (
	"api/config"
	"api/logger"
	"api/router"
	"api/utils"
	"flag"
	"net/http"
	"time"
)

func main() {
	var cfn string
	// 0.从命令行获取可能的conf路径
	// api -conf="./conf/config_qa.yaml"
	// api -conf="./conf/config_online.yaml"
	flag.StringVar(&cfn, "conf", "./config/config.yaml", "指定配置文件路径")
	flag.Parse()
	// 1. 加载配置文件
	err := config.Init(cfn)
	if err != nil {
		panic(err) // 程序启动时加载配置文件失败直接退出
	}
	// 2. 加载日志
	//err = logger.Init(config.Conf.LogConfig, config.Conf.Mode)
	//if err != nil {
	//	panic(err) // 程序启动时初始化日志模块失败直接退出
	//}
	// 3. 初始化翻译
	_ = utils.InitTrans(utils.DefaultLocale)
	r := router.InitRouter()
	server := &http.Server{
		Addr:           config.Conf.Addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("启动失败...", err)
	}
}
