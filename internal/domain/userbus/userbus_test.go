package userbus_test

import (
	"bytes"
	"context"
	"net/mail"
	"testing"
	"time"

	"github.com/hamidoujand/sales/internal/dbtest"
	"github.com/hamidoujand/sales/internal/domain/userbus"
	"github.com/hamidoujand/sales/internal/domain/userbus/userdb"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	database := dbtest.NewDatabase(ctx, t, "create_user")

	store := userdb.NewStore(database.DB)
	bus := userbus.New(store)

	email, err := mail.ParseAddress("john@gmail.com")
	if err != nil {
		t.Fatalf("parsing email: %s", err)
	}

	nu := userbus.NewUser{
		Name:     "John",
		Email:    *email,
		Roles:    []userbus.Role{userbus.RoleUser},
		Password: "password",
	}

	user, err := bus.Create(ctx, nu)
	if err != nil {
		t.Fatalf("creating user failed: %s", err)
	}

	if user.Name != nu.Name {
		t.Errorf("name=%s, got=%s", nu.Name, user.Name)
	}

	if user.Email.Address != nu.Email.Address {
		t.Errorf("email=%s, got=%s", nu.Email.Address, user.Email.Address)
	}

	if user.Enabled != true {
		t.Errorf("enabled=%t, got=%t", true, user.Enabled)
	}

	if bytes.Equal(user.PasswordHash, []byte(nu.Password)) {
		t.Errorf("password must be hashed")
	}

}
