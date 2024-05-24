package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

const (
	USER_ROLE  = "user"
	ADMIN_ROLE = "admin"
)
