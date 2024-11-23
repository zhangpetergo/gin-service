package tests

import (
	"fmt"
	"github.com/zhangpetergo/gin-service/api/http/api/apitest"
	"github.com/zhangpetergo/gin-service/app/api/page"
	"github.com/zhangpetergo/gin-service/app/domain/userapp"
	"github.com/zhangpetergo/gin-service/business/domain/userbus"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
)

func userQuery200(sd apitest.SeedData) []apitest.Table {
	usrs := make([]userbus.User, 0, len(sd.Admins)+len(sd.Users))

	for _, adm := range sd.Admins {
		usrs = append(usrs, adm.User.User)
	}

	for _, usr := range sd.Users {
		usrs = append(usrs, usr.User.User)
	}

	sort.Slice(usrs, func(i, j int) bool {
		return usrs[i].ID.String() <= usrs[j].ID.String()
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/users?page=1&rows=10&orderBy=user_id,ASC&name=Name",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &page.Document[userapp.User]{},
			ExpResp: &page.Document[userapp.User]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(usrs),
				Items:       toAppUsers(usrs),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func userQueryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &userapp.User{},
			ExpResp:    toAppUserPtr(sd.Users[0].User.User),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}