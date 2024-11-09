package userdb

import (
	"database/sql"
	"fmt"
	"github.com/zhangpetergo/gin-service/business/api/sqldb/dbarray"
	"github.com/zhangpetergo/gin-service/business/domain/userbus"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	Enabled      bool           `db:"enabled"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

// toXXX 方法的命名让我们调用时更清晰，知道 从 A to B
func toDBUser(usr userbus.User) user {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return user{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department,
			Valid:  usr.Department != "",
		},
		Enabled:     usr.Enabled,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateUpdated.UTC(),
	}
}

func toBusUser(dbUsr user) (userbus.User, error) {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]userbus.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		var err error
		roles[i], err = userbus.ParseRole(value)
		if err != nil {
			return userbus.User{}, fmt.Errorf("parse role: %w", err)
		}
	}

	bus := userbus.User{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Enabled:      dbUsr.Enabled,
		Department:   dbUsr.Department.String,
		DateCreated:  dbUsr.DateCreated.In(time.Local),
		DateUpdated:  dbUsr.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusUsers(dbUsers []user) ([]userbus.User, error) {
	bus := make([]userbus.User, len(dbUsers))

	for i, dbUsr := range dbUsers {
		var err error
		bus[i], err = toBusUser(dbUsr)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
