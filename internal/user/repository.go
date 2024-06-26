package user

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByEmail(email string) (user entity.User, err error)
	GetUserByUsername(email string) (user entity.User, err error)
	RegisterUser(user *entity.User) (userId uuid.UUID, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserByEmail(email string) (user entity.User, err error) {
	err = r.db.Get(&user, "SELECT * FROM users WHERE email = $1 and role = $2", email, entity.USER_ROLE)
	return user, err
}

func (r *repository) GetUserByUsername(username string) (user entity.User, err error) {
	err = r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return user, err
}

func (r *repository) RegisterUser(user *entity.User) (userId uuid.UUID, err error) {
	query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
	err = r.db.QueryRowx(query, user.Username, user.Email, user.Password, entity.USER_ROLE).Scan(&userId)
	return userId, err
}
