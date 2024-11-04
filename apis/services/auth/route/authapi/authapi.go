// Package authapi maintains the web based api for auth access.
package authapi

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"net/http"

	"github.com/google/uuid"
)

type api struct {
	auth *auth.Auth
}

func newAPI(auth *auth.Auth) *api {
	return &api{
		auth: auth,
	}
}

func (api *api) token(c *gin.Context) {

	kid := c.Param("kid")
	if kid == "" {
		c.Error(errs.Newf(errs.FailedPrecondition, "missing kid"))
		return
	}

	ctx := c.Request.Context()
	claims := mid.GetClaims(ctx)

	tkn, err := api.auth.GenerateToken(kid, claims)
	if err != nil {
		c.Error(errs.New(errs.Internal, err))
		return
	}

	token := struct {
		Token string `json:"token"`
	}{
		Token: tkn,
	}

	c.JSON(http.StatusOK, token)
}

func (api *api) authenticate(c *gin.Context) {
	// The middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	ctx := c.Request.Context()

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		c.Error(errs.New(errs.Unauthenticated, err))
		return
	}

	resp := struct {
		UserID uuid.UUID
		Claims auth.Claims
	}{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	c.JSON(http.StatusOK, resp)
}

func (api *api) authorize(c *gin.Context) {

	ctx := c.Request.Context()

	var auth struct {
		Claims auth.Claims
		UserID uuid.UUID
		Rule   string
	}
	if err := web.Decode(c.Request, &auth); err != nil {
		c.Error(errs.New(errs.FailedPrecondition, err))
		return
	}

	if err := api.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		c.Error(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", auth.Claims.Roles, auth.Rule, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
