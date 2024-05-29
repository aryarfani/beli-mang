package order

import (
	"beli-mang/internal/helper"
	"beli-mang/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/users/estimate", middleware.Auth(), resource.create)
}

type resource struct {
	service Service
}

func (resource resource) create(c *fiber.Ctx) error {
	var req CreateEstimationRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	userId := c.Locals("user_id").(uuid.UUID)
	order := req.ToOrder(userId)
	resp, err := resource.service.Create(order)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Estimate created successfully",
		"data":    resp,
	})
}
