package user

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByUsername(email string) (entity.User, error)
	RegisterUser(user *entity.User) (uuid.UUID, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return user, err
}

func (r *repository) RegisterUser(user *entity.User) (uuid.UUID, error) {
	var userId uuid.UUID
	query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRowx(query, user.Username, user.Email, user.Password, entity.USER_ROLE).Scan(&userId)
	return userId, err
}
