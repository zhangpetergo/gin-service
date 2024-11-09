package userapi

import (
	"github.com/zhangpetergo/gin-service/app/domain/userapp"
	"net/http"
)

func parseQueryParams(r *http.Request) (userapp.QueryParams, error) {
	values := r.URL.Query()
	filter := userapp.QueryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("row"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("user_id"),
		Name:             values.Get("name"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}
	return filter, nil
}
