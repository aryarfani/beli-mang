package user

import (
	"beli-mang/internal/entity"
	"beli-mang/pkg/hash"
	"beli-mang/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	Register(req RegisterUserRequest) (string, error)
	Login(req LoginUserRequest) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(req RegisterUserRequest) (string, error) {
	_, err := s.repo.GetUserByUsername(req.Username)
	if err == nil {
		return "", fiber.NewError(fiber.StatusConflict, "username already exists")
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
	token, _ := jwt.GenerateToken(&userId)

	return token, nil
}

func (s *service) Login(req LoginUserRequest) (string, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", err
	}
	// Generate access token
	token, _ := jwt.GenerateToken(&user.ID)

	return token, nil
}
