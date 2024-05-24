package main

import (
	"beli-mang/config"
	"beli-mang/handler"
	"beli-mang/helper"
	"errors"
	"fmt"
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

	handler := handler.HandlerConfig{
		App: app,
		DB:  configs.DB,
	}
	handler.BuildHandlers()

	log.Fatal(app.Listen(":8080"))
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
			"message": message,
		})
	}

	err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
