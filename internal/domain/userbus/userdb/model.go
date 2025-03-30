package userdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hamidoujand/sales/internal/domain/userbus"
)

type postgresUser struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Roles        []string  `db:"roles"`
	PasswordHash []byte    `db:"password_hash"`
	Enabled      bool      `db:"enabled"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

func toPostgresUser(usr userbus.User) postgresUser {
	return postgresUser{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        userbus.EncodeRoles(usr.Roles),
		PasswordHash: usr.PasswordHash,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated,
		DateUpdated:  usr.DateUpdated,
	}
}
