package admin

import (
	"beli-mang/internal/entity"
	"beli-mang/pkg/hash"
	"beli-mang/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	Register(req RegisterAdminRequest) (token string, err error)
	Login(req LoginAdminRequest) (token string, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(req RegisterAdminRequest) (token string, err error) {
	_, err = s.repo.GetUserByUsername(req.Username)
	if err == nil {
		return "", fiber.NewError(fiber.StatusConflict, "username already exists")
	}

	_, err = s.repo.GetAdminByEmail(req.Email)
	if err == nil {
		return "", fiber.NewError(fiber.StatusConflict, "admin with this email already exists")
	}

	hashedPassword := hash.HashPassword(req.Password)
	userId, err := s.repo.RegisterUser(&entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return "", err
	}

	// Generate access token
	token, _ = jwt.GenerateToken(&userId)

	return token, nil
}

func (s *service) Login(req LoginAdminRequest) (token string, err error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, "username not found")
	}

	if user.Role != entity.ADMIN_ROLE {
		return "", fiber.NewError(fiber.StatusBadRequest, "user is not user")
	}

	// Generate access token
	token, _ = jwt.GenerateToken(&user.ID)

	return token, nil
}
