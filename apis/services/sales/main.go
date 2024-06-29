package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zhangpetergo/gin-service/apis/services/api/debug"
	"github.com/zhangpetergo/gin-service/apis/services/sales/mux"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "message", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration
	cfg := struct {
		Version struct {
			Build string
			Desc  string
		}
		Web struct {
			ReadTimeout        time.Duration
			WriteTimeout       time.Duration
			IdleTimeout        time.Duration
			ShutdownTimeout    time.Duration
			APIHost            string
			DebugHost          string
			CORSAllowedOrigins []string
		}
	}{
		Version: struct {
			Build string
			Desc  string
		}{build, "sales service"},
	}

	// 设置配置默认值
	viper.SetDefault("Web.ReadTimeout", "5s")
	viper.SetDefault("Web.WriteTimeout", "10s")
	viper.SetDefault("Web.IdleTimeout", "120s")
	viper.SetDefault("Web.ShutdownTimeout", "20s")
	viper.SetDefault("Web.APIHost", "0.0.0.0:3000")
	viper.SetDefault("Web.DebugHost", "0.0.0.0:3010")
	viper.SetDefault("Web.CORSAllowedOrigins", "*")

	// 设置配置文件路径和名称
	configPath := "./zarf/config"
	configName := "config"

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	// 检查文件是否存在
	if _, err := os.Stat(configPath + "/" + configName + ".yaml"); err == nil {
		// 文件存在，读取配置文件
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		// 文件不存在，使用默认配置
	} else {
		// 其他错误
		return err
	}

	// 读取配置
	err := viper.Unmarshal(&cfg)
	if err != nil {
		// 解析配置失败
		return err
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Version.Build)
	defer log.Info(ctx, "shutdown complete")

	// 打印配置
	log.Info(ctx, "startup", "config", cfg)

	expvar.NewString("build").Set(cfg.Version.Build)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(log, shutdown),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api v1 router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		// 设置超时控制
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
