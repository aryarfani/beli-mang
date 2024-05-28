package admin

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetAdminByEmail(email string) (user entity.User, err error)
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

func (r *repository) GetAdminByEmail(email string) (user entity.User, err error) {
	err = r.db.Get(&user, "SELECT * FROM users WHERE email = $1 and role = $2 ", email, entity.ADMIN_ROLE)
	return user, err
}

func (r *repository) GetUserByUsername(username string) (user entity.User, err error) {
	err = r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return user, err
}

func (r *repository) RegisterUser(user *entity.User) (userId uuid.UUID, err error) {
	query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"
	err = r.db.QueryRowx(query, user.Username, user.Email, user.Password, entity.ADMIN_ROLE).Scan(&userId)
	return userId, err
}
