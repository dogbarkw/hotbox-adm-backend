package cli

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"

	"gorm.io/gorm/logger"

	"github.com/cloudwego/kitex/pkg/klog"
	dmysql "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Dialer struct {
	client *ssh.Client
	_      *context.Context
}

func (v *Dialer) Dial(context context.Context, address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

func DialWithKeyFile() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            os.Getenv("SSH_USER"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if k, err := os.ReadFile(os.Getenv("SSH_KEYFILE")); err != nil {
		klog.Errorf("ecommerce.cli.ReadFile error:%v", err)
		return nil, err
	} else {
		signer, err := ssh.ParsePrivateKey(k)
		if err != nil {
			klog.Errorf("ecommerce.cli.ParsePrivateKey error:%v", err)
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	address := fmt.Sprintf("%s:%s", os.Getenv("SSH_HOST"), os.Getenv("SSH_PORT"))
	return ssh.Dial("tcp", address, config)
}

func NewSshDB() (db *gorm.DB, err error) {
	dial, err := DialWithKeyFile()
	if err != nil {
		return nil, err
	}
	// defer dial.Close()

	// 注册ssh代理
	// dmysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)
	dmysql.RegisterDialContext("mysql+ssh", (&Dialer{client: dial}).Dial)
	// 填写注册的mysql网络
	dsn := fmt.Sprintf("%s:%s@mysql+ssh(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DB"))
	logLevel := logger.Silent
	if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
		lvl, err := strconv.Atoi(logLevelStr)
		if err != nil {
			logLevel = logger.Silent
		} else {
			logLevel = logger.LogLevel(lvl)
		}
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSshHDTaskDB() (db *gorm.DB, err error) {
	dial, err := DialWithKeyFile()
	if err != nil {
		return nil, err
	}
	// defer dial.Close()

	// 注册ssh代理
	// dmysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)
	dmysql.RegisterDialContext("mysql+ssh", (&Dialer{client: dial}).Dial)
	// 填写注册的mysql网络
	dsn := fmt.Sprintf("%s:%s@mysql+ssh(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_TASK_USER"), os.Getenv("HD_TASK_PASS"), os.Getenv("HD_TASK_HOST"), os.Getenv("HD_TASK_PORT"), os.Getenv("HD_TASK_DB"))
	logLevel := logger.Silent
	if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
		lvl, err := strconv.Atoi(logLevelStr)
		if err != nil {
			logLevel = logger.Silent
		} else {
			logLevel = logger.LogLevel(lvl)
		}
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewHDSshDB() (db *gorm.DB, err error) {
	dial, err := DialWithKeyFile()
	if err != nil {
		return nil, err
	}
	// defer dial.Close()

	// 注册ssh代理
	// dmysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)
	dmysql.RegisterDialContext("mysql+ssh", (&Dialer{client: dial}).Dial)
	// 填写注册的mysql网络
	dsn := fmt.Sprintf("%s:%s@mysql+ssh(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_MYSQL_USER"), os.Getenv("HD_MYSQL_PASS"), os.Getenv("HD_MYSQL_HOST"), os.Getenv("HD_MYSQL_PORT"), os.Getenv("HD_MYSQL_DB"))
	logLevel := logger.Silent
	if logLevelStr := os.Getenv("MYSQL_LOG_LEVEL"); logLevelStr != "" {
		lvl, err := strconv.Atoi(logLevelStr)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi error: %w", err)
		}

		logLevel = logger.LogLevel(lvl)
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewHdADBSshDB() (db *sql.DB, err error) {
	dial, err := DialWithKeyFile()
	if err != nil {
		return nil, err
	}
	// defer dial.Close()

	// 注册ssh代理
	// dmysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)
	dmysql.RegisterDialContext("mysql+ssh", (&Dialer{client: dial}).Dial)
	// 填写注册的mysql网络
	dsn := fmt.Sprintf("%s:%s@mysql+ssh(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("HD_ADB_USER"), os.Getenv("HD_ADB_PASS"), os.Getenv("HD_ADB_HOST"), os.Getenv("HD_ADB_PORT"), os.Getenv("HD_ADB_DB"))
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
