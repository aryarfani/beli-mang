package handler

import (
	"beli-mang/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerConfig struct {
	App *fiber.App
	DB  *sqlx.DB
}

func (handler HandlerConfig) BuildHandlers() {
	user.RegisterHandlers(handler.App, user.NewService(user.NewRepository(handler.DB)))
}
