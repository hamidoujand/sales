package userbus

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/hamidoujand/sales/internal/order"
	"github.com/hamidoujand/sales/internal/page"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// storer represents the required behavior from the storage engine.
type storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
}

type UserBus struct {
	store storer
}

func New(store storer) *UserBus {
	return &UserBus{
		store: store,
	}
}

func (u *UserBus) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hashing password: %w", err)
	}

	now := time.Now()
	usr := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		PasswordHash: hash,
		Email:        nu.Email,
		Roles:        nu.Roles,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := u.store.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("creating user: %w", err)
	}
	return usr, nil
}

func (u *UserBus) Update(ctx context.Context, usr User, updates UpdateUser) (User, error) {
	if updates.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*updates.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("hashing password: %w", err)
		}
		usr.PasswordHash = hash
	}

	if updates.Email != nil {
		usr.Email = *updates.Email
	}

	if updates.Roles != nil {
		usr.Roles = updates.Roles
	}

	if updates.Name != nil {
		usr.Name = *updates.Name
	}

	if updates.Enabled != nil {
		usr.Enabled = *updates.Enabled
	}

	usr.DateUpdated = time.Now()

	if err := u.store.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("updating user: %w", err)
	}

	return usr, nil
}

func (u *UserBus) Delete(ctx context.Context, usr User) error {
	if err := u.store.Delete(ctx, usr); err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return nil
}

func (u *UserBus) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	usr, err := u.store.QueryByID(ctx, userID)
	if err != nil {
		//check for not-found
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("query by ID: %w", err)
	}

	return usr, nil
}

func (u *UserBus) Query(ctx context.Context, filter QueryFilter, order order.By, page page.Page) ([]User, error) {
	users, err := u.store.Query(ctx, filter, order, page)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	return users, nil
}
