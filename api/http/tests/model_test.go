package tests

import (
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/app/domain/userapp"
	"github.com/zhangpetergo/gin-service/business/domain/userbus"
	"time"
)

func toErrorPtr(err errs.Error) *errs.Error {
	return &err
}

func toAppUser(usr userbus.User) userapp.User {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return userapp.User{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: nil, // This field is not marshalled.
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []userapp.User {
	items := make([]userapp.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(usr userbus.User) *userapp.User {
	appUsr := toAppUser(usr)
	return &appUsr
}
