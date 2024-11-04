package mid

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// AuthorizeService executes the specified role and does not extract any domain data.
func AuthorizeService(log *logger.Logger, client *authclient.Client, rule string) gin.HandlerFunc {

	return func(c *gin.Context) {

		if len(c.Errors) > 0 {
			// 如果有错，直接退出
			return
		}

		ctx := c.Request.Context()

		userID, err := GetUserID(ctx)
		if err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
		}

		auth := authclient.Authorize{
			Claims: GetClaims(ctx),
			UserID: userID,
			Rule:   rule,
		}

		if err := client.Authorize(ctx, auth); err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
		}

		c.Next()
	}

}

// AuthorizeLocal executes the specified role and does not extract any domain data.
func AuthorizeLocal(auth *auth.Auth, rule string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID, err := GetUserID(ctx)

		if err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
		}

		if err := auth.Authorize(ctx, GetClaims(ctx), userID, rule); err != nil {
			c.Error(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", GetClaims(ctx).Roles, rule, err))
		}
		c.Next()
	}
}
