package userbus

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
)

// QueryFilter represents all the fields that can be used for filtering.
type QueryFilter struct {
	ID             *uuid.UUID
	Name           *string
	Email          *mail.Address
	StartCreatedAt *time.Time
	EndCreatedAt   *time.Time
}
