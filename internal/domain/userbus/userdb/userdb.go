package userdb

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/hamidoujand/sales/internal/domain/userbus"
	"github.com/hamidoujand/sales/internal/order"
	"github.com/hamidoujand/sales/internal/page"
	"github.com/hamidoujand/sales/internal/sqldb"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

// Delete implements userbus.storer.
func (s *Store) Delete(ctx context.Context, usr userbus.User) error {
	panic("unimplemented")
}

// Query implements userbus.storer.
func (s *Store) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, page page.Page) ([]userbus.User, error) {
	panic("unimplemented")
}

// QueryByEmail implements userbus.storer.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	panic("unimplemented")
}

// QueryByID implements userbus.storer.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	panic("unimplemented")
}

// Update implements userbus.storer.
func (s *Store) Update(ctx context.Context, usr userbus.User) error {
	panic("unimplemented")
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, usr userbus.User) error {
	const q = `
	INSERT INTO users(id,name,email,password_hash,roles,enabled,date_created,date_updated)
	VALUES (:id,:name,:email,:password_hash,:roles,:enabled,:date_created,:date_updated);
	`
	if err := sqldb.NamedExecContext(ctx, s.db, q, toPostgresUser(usr)); err != nil {
		if errors.Is(err, sqldb.ErrDuplicatedEntry) {
			return userbus.ErrDuplicatedEmail
		} else {
			return fmt.Errorf("namedExecContext: %w", err)
		}
	}
	return nil
}
