package mid

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/business/domain/userbus"
	"github.com/zhangpetergo/gin-service/foundation/logger"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize executes the specified role and does not extract any domain data.
func Authorize(log *logger.Logger, client *authclient.Client, rule string) gin.HandlerFunc {

	return func(c *gin.Context) {

		if len(c.Errors) > 0 {
			// 如果有错，直接退出
			return
		}

		ctx := c.Request.Context()

		// 在使用日志调试的过程中，发现需要一条全局的日志 logger 还是很有必要的
		// 是不是可以更改 logger 的操作
		log.Info(ctx, "-----------------Authorize------------------------")

		userID, err := GetUserID(ctx)
		if err != nil {
			log.Info(ctx, "-----------------FUCK-----------------")
			c.Error(errs.New(errs.Unauthenticated, err))
			return
		}

		auth := authclient.Authorize{
			Claims: GetClaims(ctx),
			UserID: userID,
			Rule:   rule,
		}

		if err := client.Authorize(ctx, auth); err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
			return
		}

		c.Next()
	}

}

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(log *logger.Logger, client *authclient.Client, userBus *userbus.Business, rule string) gin.HandlerFunc {

	return func(c *gin.Context) {

		if len(c.Errors) > 0 {
			// 如果有错，直接退出
			return
		}

		ctx := c.Request.Context()

		id := c.Param("user_id")

		var userID uuid.UUID

		if id != "" {
			var err error
			userID, err = uuid.Parse(id)
			if err != nil {
				c.Error(errs.New(errs.Unauthenticated, ErrInvalidID))
				return
			}

			usr, err := userBus.QueryByID(ctx, userID)
			if err != nil {
				switch {
				case errors.Is(err, userbus.ErrNotFound):
					c.Error(errs.New(errs.Unauthenticated, err))
				default:
					c.Error(errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err))
				}
			}

			ctx = setUser(ctx, usr)
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

//// AuthorizeLocal executes the specified role and does not extract any domain data.
//func AuthorizeLocal(auth *auth.Auth, rule string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx := c.Request.Context()
//
//		userID, err := GetUserID(ctx)
//
//		if err != nil {
//			c.Error(errs.New(errs.Unauthenticated, err))
//		}
//
//		if err := auth.Authorize(ctx, GetClaims(ctx), userID, rule); err != nil {
//			c.Error(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", GetClaims(ctx).Roles, rule, err))
//		}
//		c.Next()
//	}
//}
