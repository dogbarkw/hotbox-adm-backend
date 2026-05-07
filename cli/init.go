package cli

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var (
	RedisCli         redis.Cmdable
	GormDB           *gorm.DB
	HotDogGormDB     *gorm.DB
	HotDogTaskGormDB *gorm.DB
	HotDogRedis      redis.Cmdable
	SpecialUserIds   []string
	HotDogADBGormDB  *sql.DB
)

func InitEnv() {
	// 读取配置文件，没有的话就读本机的环境变量
	_ = godotenv.Load()
	paths := []string{
		"./.env",
		"../.env",
		"../../.env",
		"../../../.env",
	}
	var e error
	for _, v := range paths {
		err := godotenv.Load(v)
		e = err
		if err == nil {
			break
		}
	}
	if e != nil {
		panic(e)
	}
}

func InitSpecialUserIds() {
	var userIds []int64
	err := HotDogGormDB.WithContext(context.Background()).Table("test_user").Where("status", 1).Select("user_id").Scan(&userIds).Error
	if err != nil {
		panic(err)
	}

	for _, userId := range userIds {
		SpecialUserIds = append(SpecialUserIds, cast.ToString(userId))
	}
	// 避免SpecialUserIds MySQL查询的时候in NULL
	if len(SpecialUserIds) == 0 {
		SpecialUserIds = append(SpecialUserIds, "0")
	}
	logrus.Infof("SpecialUserIds:%v", SpecialUserIds)
}

func InitGormDB() {
	var err error
	// if strings.ToLower(os.Getenv("ENV")) == "dev" && os.Getenv("MYSQL_HOST") != "localhost" {
	GormDB, err = NewSshDB()
	// } else {
	// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DB"))
	// 	logLevel := logger.Silent
	// 	if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
	// 		lvl, err := strconv.Atoi(logLevelStr)
	// 		if err != nil {
	// 			logLevel = logger.Silent
	// 		} else {
	// 			logLevel = logger.LogLevel(lvl)
	// 		}
	// 	}
	// 	GormDB, err = gorm.Open(mysql.Open(dsn),
	// 		&gorm.Config{
	// 			PrepareStmt:            true,
	// 			SkipDefaultTransaction: true,
	// 			NamingStrategy: schema.NamingStrategy{
	// 				SingularTable: true, // 使用单数表名
	// 			},
	// 			Logger: logger.Default.LogMode(logLevel),
	// 		},
	// 	)
	// }
	if err != nil {
		panic(err)
	}

	sdb, err := GormDB.DB()
	if err != nil {
		panic(err)
	}
	sdb.SetMaxIdleConns(10)                  // 最大空闲连接数
	sdb.SetMaxOpenConns(100)                 // 最大连接数
	sdb.SetConnMaxLifetime(time.Minute * 10) // 设置连接空闲超时
}

func InitHDTaskDB() {
	var err error
	if strings.ToLower(os.Getenv("ENV")) == "dev" && os.Getenv("MYSQL_HOST") != "localhost" {
		HotDogTaskGormDB, err = NewSshHDTaskDB()
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_TASK_USER"), os.Getenv("HD_TASK_PASS"), os.Getenv("HD_TASK_HOST"), os.Getenv("HD_TASK_PORT"), os.Getenv("HD_TASK_DB"))
		logLevel := logger.Silent
		if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
			lvl, err := strconv.Atoi(logLevelStr)
			if err != nil {
				logLevel = logger.Silent
			} else {
				logLevel = logger.LogLevel(lvl)
			}
		}
		HotDogTaskGormDB, err = gorm.Open(mysql.Open(dsn),
			&gorm.Config{
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true, // 使用单数表名
				},
				Logger: logger.Default.LogMode(logLevel),
			},
		)
	}
	if err != nil {
		panic(err)
	}

	sdb, err := HotDogTaskGormDB.DB()
	if err != nil {
		panic(err)
	}
	sdb.SetMaxIdleConns(10)                  // 最大空闲连接数
	sdb.SetMaxOpenConns(100)                 // 最大连接数
	sdb.SetConnMaxLifetime(time.Minute * 10) // 设置连接空闲超时
}

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		DialTimeout:  time.Second * 2,
		PoolSize:     50,
		MinIdleConns: 30,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	RedisCli = client
}

func InitHDGormDB() {
	var err error
	if os.Getenv("ENV") == "dev" {
		HotDogGormDB, err = NewHDSshDB()
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_MYSQL_USER"), os.Getenv("HD_MYSQL_PASS"), os.Getenv("HD_MYSQL_HOST"), os.Getenv("HD_MYSQL_PORT"), os.Getenv("HD_MYSQL_DB"))
		logLevel := logger.Silent
		if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
			lvl, err := strconv.Atoi(logLevelStr)
			if err != nil {
				logLevel = logger.Silent
			} else {
				logLevel = logger.LogLevel(lvl)
			}
		}
		HotDogGormDB, err = gorm.Open(mysql.Open(dsn),
			&gorm.Config{
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true, // 使用单数表名
				},
				Logger: logger.Default.LogMode(logLevel),
			},
		)
	}

	if err != nil {
		panic(err)
	}

	sdb, err := HotDogGormDB.DB()
	if err != nil {
		panic(err)
	}
	sdb.SetMaxIdleConns(10)                  // 最大空闲连接数
	sdb.SetMaxOpenConns(100)                 // 最大连接数
	sdb.SetConnMaxLifetime(time.Minute * 10) // 设置连接空闲超时
}

func InitHDRedis() {
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HD_REDIS_HOST"), os.Getenv("HD_REDIS_PORT")),
		ReadTimeout:  -2,
		WriteTimeout: -2,
		DialTimeout:  time.Second * 2,
		PoolSize:     50,
		MinIdleConns: 30,
	}
	if os.Getenv("ENV") == "dev" {
		dial, err := DialWithKeyFile()
		if err != nil {
			panic(err)
		}
		opt.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dial.Dial(network, addr)
		}
	}

	client := redis.NewClient(opt)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	HotDogRedis = client
}

func InitHDADBGormDB() {
	var err error
	HotDogADBGormDB, err = NewHdADBSshDB()
	// if strings.ToLower(os.Getenv("ENV")) == "dev" && os.Getenv("HD_ADB_HOST") != "localhost" {
	// } else {
	// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_ADB_USER"), os.Getenv("HD_ADB_PASS"), os.Getenv("HD_ADB_HOST"), os.Getenv("HD_ADB_PORT"), os.Getenv("HD_ADB_DB"))
	// 	HotDogADBGormDB, err = sql.Open("mysql", dsn)
	// }
	if err != nil {
		panic(err)
	}
	sdb := HotDogADBGormDB
	sdb.SetMaxIdleConns(10)                  // 最大空闲连接数
	sdb.SetMaxOpenConns(100)                 // 最大连接数
	sdb.SetConnMaxLifetime(time.Minute * 10) // 设置连接空闲超时
}
