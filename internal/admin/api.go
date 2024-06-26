package admin

import (
	"beli-mang/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/admin/register", resource.register)
	app.Post("/admin/login", resource.login)
}

type resource struct {
	service Service
}

func (resource resource) register(c *fiber.Ctx) error {
	var req RegisterAdminRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	token, err := resource.service.Register(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin registered successfully",
		"token":   token,
	})
}

func (resource resource) login(c *fiber.Ctx) error {
	var req LoginAdminRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	token, err := resource.service.Login(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin login successfully",
		"token":   token,
	})
}
