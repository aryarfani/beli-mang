package handler

import (
	"beli-mang/internal/admin"
	"beli-mang/internal/image"
	"beli-mang/internal/merchant"
	"beli-mang/internal/user"
	"beli-mang/pkg/aws"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerConfig struct {
	App      *fiber.App
	DB       *sqlx.DB
	S3Config *aws.S3Config
}

func (handler HandlerConfig) BuildHandlers() {
	user.RegisterHandlers(handler.App, user.NewService(user.NewRepository(handler.DB)))
	admin.RegisterHandlers(handler.App, admin.NewService(admin.NewRepository(handler.DB)))
	image.RegisterHandlers(handler.App, image.NewService(handler.S3Config))
	merchant.RegisterHandlers(handler.App, merchant.NewService(merchant.NewRepository(handler.DB)))
}
