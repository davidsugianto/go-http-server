package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidsugianto/go-http-server/internal/pkg/config"
	"github.com/davidsugianto/go-pkgs/logs"
)

var (
	errLogPath   string
	infoLogPath  string
	debugLogPath string
)

func main() {
	flag.StringVar(&errLogPath, "error_log", "/tmp/idp-core.error.log", "error log")
	flag.StringVar(&infoLogPath, "info_log", "/tmp/idp-core.info.log", "info log")
	flag.StringVar(&debugLogPath, "debug_log", "/tmp/idp-core.debug.log", "debug log")
	flag.Parse()

	setLog(logs.ErrorLevel, errLogPath)
	setLog(logs.InfoLevel, infoLogPath)
	setLog(logs.DebugLevel, debugLogPath)

	logs.Info("Starting Go HTTP Server in %s mode", os.Getenv("ENV"))

	// load config
	cfgPath := fmt.Sprintf("configs/config.%s.yaml", os.Getenv("ENV"))
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logs.Fatalf("cannot load config from %s", cfgPath)
		panic(err)
	}

	// init client

	// init repository

	// init use case

	// init server dependencies
	server := New(Dependencies{
		Config: cfg,
	})

	// run server
	logs.Info("listening on port")
	if err := server.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		logs.Fatal("server failed to run")
	}
}

func setLog(level logs.Level, filePath string) {
	lgr, err := logs.NewLogger(&logs.Config{
		Level:   level,
		LogFile: filePath,
		Caller:  false,
		AppName: "go-http-server",
		UseJSON: false,
	})
	if err != nil {
		logs.Fatalln(err)
	}

	err = logs.SetLogger(level, lgr)
	if err != nil {
		logs.Fatalln(err)
	}
}
