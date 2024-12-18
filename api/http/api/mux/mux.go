// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"context"
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build      string
	Log        *logger.Logger
	Auth       *auth.Auth
	AuthClient *authclient.Client
	DB         *sqlx.DB
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) *web.App {
	logger := func(ctx context.Context, msg string, v ...any) {
		cfg.Log.Info(ctx, msg, v...)
	}

	app := web.NewApp(
		logger,
	)

	routeAdder.Add(app, cfg)

	return app
}
