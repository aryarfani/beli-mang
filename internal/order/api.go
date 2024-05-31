package order

import (
	"beli-mang/internal/helper"
	"beli-mang/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/users/estimate", middleware.Auth(), resource.createEstimate)
	app.Post("/users/orders", middleware.Auth(), resource.createOrder)
	app.Get("/users/orders", middleware.Auth(), resource.query)
}

type resource struct {
	service Service
}

func (resource resource) createEstimate(c *fiber.Ctx) error {
	var req CreateEstimationRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	userId := c.Locals("user_id").(uuid.UUID)
	orders := req.ToOrders(userId)
	resp, err := resource.service.Create(orders, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (resource resource) createOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	resp, err := resource.service.CreateOrder(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (resource resource) query(c *fiber.Ctx) error {
	var params QueryOrdersRequest
	_ = c.QueryParser(&params)

	params.UserId = c.Locals("user_id").(uuid.UUID)
	items, err := resource.service.Query(params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(items)
}
