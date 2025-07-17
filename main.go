package main

import (
	"blog/config"
	"blog/middleware"
	"blog/routes"
	"blog/storage/mongodb"
	"blog/storage/mysql"
	"blog/storage/redis"
	"blog/utils"
	"log"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer utils.Log.Sync()

	// 加载配置
	if err := config.Init(); err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// 初始化数据库连接
	if err := mysql.Init(); err != nil {
		utils.Fatal("初始化MySQL失败", zap.Error(err))
	}

	// 初始化MongoDB连接
	if err := mongodb.Init(); err != nil {
		utils.Fatal("初始化MongoDB失败", zap.Error(err))
	}

	// 初始化Redis连接
	if err := redis.Init(); err != nil {
		utils.Fatal("初始化Redis失败", zap.Error(err))
	}

	// 初始化路由
	r := routes.Init()

	// 使用日志中间件
	r.Use(middleware.Logger())

	// 启动服务器
	port := config.GlobalConfig.Server.Port
	utils.Info("服务器启动", zap.Int("port", port))
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		utils.Fatal("服务器启动失败", zap.Error(err))
	}
}
