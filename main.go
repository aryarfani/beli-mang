package main

import (
	"beli-mang/config"

	"beli-mang/internal/admin"
	"beli-mang/internal/helper"
	"beli-mang/internal/image"
	"beli-mang/internal/item"
	"beli-mang/internal/merchant"
	"beli-mang/internal/order"
	"beli-mang/internal/user"
	"errors"
	"log"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	configs, err := config.LoadConfigs()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configs.DB.Close()

	helper.InitValidator()

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "Beli Mang",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: appErrorHandler,
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())

	registerHandlers(app, configs)

	log.Fatal(app.Listen(":8080"))
}

func registerHandlers(app *fiber.App, configs config.Configs) {
	user.RegisterHandlers(app, user.NewService(user.NewRepository(configs.DB)))
	admin.RegisterHandlers(app, admin.NewService(admin.NewRepository(configs.DB)))
	image.RegisterHandlers(app, image.NewService(&configs.S3Config))
	merchant.RegisterHandlers(app, merchant.NewService(merchant.NewRepository(configs.DB)))
	item.RegisterHandlers(app, item.NewService(item.NewRepository(configs.DB)))
	order.RegisterHandlers(app, order.NewService(order.NewRepository(configs.DB)))
}

func appErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	var message string
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	// Set default error message as json
	if message != "" {
		return ctx.Status(code).JSON(fiber.Map{
			"error": message,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
