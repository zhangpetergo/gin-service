package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zhangpetergo/gin-service/api/cmd/services/auth/build/all"
	"github.com/zhangpetergo/gin-service/api/http/api/debug"
	"github.com/zhangpetergo/gin-service/api/http/api/mux"
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/business/api/sqldb"
	"github.com/zhangpetergo/gin-service/foundation/keystore"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
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
		return web.GetTraceID(ctx)
	}
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", traceIDFn, events)
	// -------------------------------------------------------------------------
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
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
		Auth struct {
			KeysFolder string
			ActiveKID  string
			Issuer     string
		}
		DB struct {
			User         string
			Password     string
			HostPort     string
			Name         string
			MaxIdleConns int
			MaxOpenConns int
			DisableTLS   bool
		}
	}{
		Version: struct {
			Build string
			Desc  string
		}{build, "sales service"},
	}
	// 设置配置默认值
	// web
	viper.SetDefault("Web.ReadTimeout", "5s")
	viper.SetDefault("Web.WriteTimeout", "10s")
	viper.SetDefault("Web.IdleTimeout", "120s")
	viper.SetDefault("Web.ShutdownTimeout", "20s")
	viper.SetDefault("Web.APIHost", "0.0.0.0:6000")
	viper.SetDefault("Web.DebugHost", "0.0.0.0:6100")
	viper.SetDefault("Web.CORSAllowedOrigins", "*")

	// auth
	viper.SetDefault("Auth.KeysFolder", "zarf/keys/")
	viper.SetDefault("Auth.ActiveKID", "default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1")
	viper.SetDefault("Auth.Issuer", "default:service project")

	// DB
	viper.SetDefault("DB.User", "postgres")
	viper.SetDefault("DB.Password", "postgres")
	viper.SetDefault("DB.HostPort", "database-service.sales-system.svc.cluster.local")
	viper.SetDefault("DB.Name", "postgres")
	viper.SetDefault("DB.MaxIdleConns", "2")
	viper.SetDefault("DB.MaxOpenConns", "0")
	viper.SetDefault("DB.DisableTLS", true)

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
	// Database Support
	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.HostPort)
	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		HostPort:     cfg.DB.HostPort,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer db.Close()

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		KeyLookup: ks,
	}

	ath, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

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

	cfgMux := mux.Config{
		Build: build,
		Log:   log,
		Auth:  ath,
		DB:    db,
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfgMux, all.Routes()),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
