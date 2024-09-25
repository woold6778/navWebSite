package mydb

import (
	"context"
	"database/sql"
	"fmt"
	"nav-web-site/config"
	"nav-web-site/util"
	"time"

	"github.com/go-redis/redis/v8"
	// 添加 MySQL 驱动
	_ "github.com/go-sql-driver/mysql"
)

type TABLES struct {
	Nav        StructNav
	NavClass   StructNavClass
	News       StructNews
	NewsClass  StructNewsClass
	Admin      StructAdmin
	UploadFile StructUploadFile
	// 其他表如 User, Product 等都可以类似嵌入
}

var (
	Db          *sql.DB
	RedisClient *redis.Client // 全局 Redis 客户端
	Ctx         = context.Background()
	Tables      TABLES // 全局 TABLES 实例
)

func InitDB() {
	var err error

	fmt.Println("Nav and Websocket Server starting...")

	// 建立数据库连接
	Db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Config.MySQL.Username,
		config.Config.MySQL.Password,
		config.Config.MySQL.Host,
		config.Config.MySQL.Port,
		config.Config.MySQL.Database,
	))
	if err != nil {
		util.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	util.InfoLogger.Println("Database connection successful")

	if Db == nil {
		util.ErrorLogger.Println("Database connection is nil")
	}

	// 设置数据库连接池参数
	Db.SetMaxOpenConns(config.Config.MySQL.MaxOpenConns)                                    // 设置最大打开连接数
	Db.SetMaxIdleConns(config.Config.MySQL.MaxIdleConns)                                    // 设置最大空闲连接数
	Db.SetConnMaxLifetime(time.Duration(config.Config.MySQL.ConnMaxLifetime) * time.Minute) // 设置连接的最大生命周期

	//连接 Redis:
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	// 测试redis连接
	_, redis_err := RedisClient.Ping(Ctx).Result()
	if redis_err != nil {
		util.ErrorLogger.Fatalf("Could not connect to Redis: %v", redis_err)
	}
}
