package merchant

import (
	"beli-mang/internal/helper"
	"beli-mang/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/admins/merchants", middleware.Auth(), resource.create)
	app.Get("/admins/merchants", middleware.Auth(), resource.query)
}

type resource struct {
	service Service
}

func (resource resource) create(c *fiber.Ctx) error {
	var req CreateMerchantRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	token, err := resource.service.Create(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Merchant created successfully",
		"token":   token,
	})
}

func (resource resource) query(c *fiber.Ctx) error {
	merchants, err := resource.service.Query()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(merchants)
}
