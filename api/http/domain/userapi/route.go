package userapi

import (
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/app/domain/userapp"
	"github.com/zhangpetergo/gin-service/business/domain/userbus"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	UserBus    *userbus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	authen := mid.Authenticate(cfg.Log, cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.Log, cfg.AuthClient, auth.RuleAdminOnly)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.Log, cfg.AuthClient, cfg.UserBus, auth.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.Log, cfg.AuthClient, cfg.UserBus, auth.RuleAdminOnly)
	api := newAPI(userapp.NewApp(cfg.UserBus))

	app.Use(mid.Trace(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
	app.GET("/users", authen, ruleAdmin, mid.CheckError(), api.query)
	app.GET("/users/:user_id", ruleAdmin, authen, mid.CheckError(), api.queryByID)
	app.POST("/users", ruleAdmin, authen, mid.CheckError(), api.create)
	app.PUT("/users/role/:user_id", ruleAuthorizeAdmin, authen, mid.CheckError(), api.updateRole)
	app.PUT("/users/:user_id", ruleAuthorizeUser, authen, mid.CheckError(), api.update)
	app.DELETE("/users/:user_id", ruleAdmin, authen, mid.CheckError(), api.delete)

}
