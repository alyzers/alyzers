package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/alyzers/alyzers/internal/alyzers/conf"
	"github.com/alyzers/alyzers/internal/alyzers/router"
	"github.com/alyzers/alyzers/pkg/cache"
	"github.com/alyzers/alyzers/pkg/ctx"
	"github.com/alyzers/alyzers/pkg/database"
	httpx "github.com/alyzers/alyzers/pkg/http"
	"github.com/alyzers/alyzers/pkg/log"
	"github.com/alyzers/alyzers/pkg/runner"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/9/4 19:51
 * @file: main.go
 * @description: alyzers program
 */

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "conf", "conf.d/config.toml", "conf file path, e.g. -conf ./conf.d")
}

func main() {
	flag.Parse()
	printRunner()

	var appConf conf.AppConfig
	appConf = conf.NewConf(configFile)

	logger := log.NewLog(&appConf.Log)

	redis, err := cache.NewRedis(appConf.Redis)
	if err != nil {
		panic(err)
	}

	// db
	db, err := database.NewDatabase(appConf.Database, *logger)
	if err != nil {
		panic(err)
	}
	mongo, err := database.NewMongoDB(appConf.Database.MongoDB, context.Background())
	if err != nil {
		panic(err)
	}
	Ctx := ctx.NewContext(context.Background(), mongo, redis, db, logger.Sugar())

	route := router.NewRouter(&appConf.Http, Ctx)
	// http srv
	cleanup := httpx.NewHttp(appConf.Http, route.Router(logger))
	cleanup()
}

func printRunner() {
	fmt.Println("runner.pwd:", runner.Pwd)
	fmt.Println("runner.hostname:", runner.Hostname)
}
