package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"

	"alertstore/internal"
	"alertstore/internal/db"
	"alertstore/internal/server"
	"alertstore/version"
)

// Args 定义 alertstore 传入参数的结构体变量
type Args struct {
	Address                string
	DBBackend              string
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int

	Debug   bool
	DryRun  bool
	Version bool
}

func main() {
	args := Args{}
	// 获取程序版本
	flag.BoolVar(&args.Version, "version", false, "print the version and exit")
	// 设置程序绑定端口, 默认为 9567
	flag.StringVar(&args.Address, "listen.address", envWithDefault("ALERTSTORE_ADDR", ":9567"), "address in which to listen for http requests, (also ALERTSTORE_ADDR)")
	// 配置是否开启 debug 级别, 是否打开 Logrus Debug 级别, 默认为 false
	flag.BoolVar(&args.Debug, "debug", false, "enable debug mode, which dumps alerts payloads to the log as they arrive")
	// 设置后端数据存储, 支持 Mysql 与 postgresql.
	flag.StringVar(&args.DBBackend, "database-backend", envWithDefault("ALERTSTORE_BACKEND", "mysql"), "database backend, allowed are mysql, postgres, and null (also ALERTSTORE_BACKEND")
	// 设置后端数据存储的 DSN. Mysql-如: ${MYSQL_USER}:${MYSQL_PASSWORD}@/${MYSQL_DATABASE}
	flag.StringVar(&args.DSN, "dsn", os.Getenv(internal.DSNVar), "Database DSN (also ALERTSTORE_DSN)")
	// 设置数据库最大连接数.
	flag.IntVar(&args.MaxOpenConns, "max-open-connections", 2, "maximum number of connections in the pool")
	// 设置数据库最大闲置连接数.
	flag.IntVar(&args.MaxIdleConns, "max-idle-connections", 1, "maximum number of idle connections in the pool")
	// 设置数据库连接使用的最大时长. 默认为 600s
	flag.IntVar(&args.MaxConnLifetimeSeconds, "max-connection-lifetyme-seconds", 600, "maximum number of seconds a connection is kept alive in the pool")

	flag.Parse()

	if args.Version {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if args.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	driver, err := db.Connect(args.DBBackend, db.ConnectionArgs{
		DSN:                    args.DSN,
		MaxIdleConns:           args.MaxIdleConns,
		MaxOpenConns:           args.MaxOpenConns,
		MaxConnLifetimeSeconds: args.MaxConnLifetimeSeconds,
	})
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		os.Exit(1)
	}

	s := server.New(driver, args.Debug)
	s.Start(args.Address)
}

// 获取系统传入环境变量
func envWithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
