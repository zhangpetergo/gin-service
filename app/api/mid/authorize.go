package mid

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/business/api/auth"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize executes the specified role and does not extract any domain data.
func Authorize(auth *auth.Auth, rule string) gin.HandlerFunc {
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
