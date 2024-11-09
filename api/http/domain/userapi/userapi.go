// Package userapi maintains the web based api for user access.
package userapi

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/app/domain/userapp"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"net/http"
)

type api struct {
	userApp *userapp.App
}

func newAPI(userApp *userapp.App) *api {
	return &api{
		userApp: userApp,
	}
}

func (api *api) create(c *gin.Context) {
	var app userapp.NewUser

	r := c.Request
	ctx := r.Context()

	if err := web.Decode(r, &app); err != nil {
		c.Error(errs.New(errs.FailedPrecondition, err))
		return
	}
	usr, err := api.userApp.Create(ctx, app)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, usr)
}
func (api *api) update(c *gin.Context) {
	var app userapp.UpdateUser

	r := c.Request
	ctx := r.Context()

	if err := web.Decode(r, &app); err != nil {
		c.Error(errs.New(errs.FailedPrecondition, err))
		return
	}
	usr, err := api.userApp.Update(ctx, app)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, usr)
}
func (api *api) updateRole(c *gin.Context) {
	var app userapp.UpdateUserRole

	r := c.Request
	ctx := r.Context()

	if err := web.Decode(r, &app); err != nil {
		c.Error(errs.New(errs.FailedPrecondition, err))
		return
	}
	usr, err := api.userApp.UpdateRole(ctx, app)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, usr)
}
func (api *api) delete(c *gin.Context) {
	ctx := c.Request.Context()

	if err := api.userApp.Delete(ctx); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
func (api *api) query(c *gin.Context) {

	r := c.Request
	ctx := r.Context()

	qp, err := parseQueryParams(r)
	if err != nil {
		c.Error(err)
		return
	}
	usr, err := api.userApp.Query(ctx, qp)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, usr)
}
func (api *api) queryByID(c *gin.Context) {

	ctx := c.Request.Context()

	usr, err := api.userApp.QueryByID(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, usr)
}
