package item

import (
	"beli-mang/internal/helper"
	"beli-mang/internal/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/admin/merchants/:merchantId/items", middleware.Auth(), resource.create)
	app.Get("/admin/merchants/:merchantId/items", middleware.Auth(), resource.query)
}

type resource struct {
	service Service
}

func (resource resource) create(c *fiber.Ctx) error {
	var req CreateItemRequest
	_ = c.BodyParser(&req)
	if err := helper.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	req.MerchantId = c.Params("merchantId")
	fmt.Println(c.Params("merchantId"))

	token, err := resource.service.Create(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Item created successfully",
		"token":   token,
	})
}

func (resource resource) query(c *fiber.Ctx) error {
	var params QueryItemsRequest
	_ = c.QueryParser(&params)

	params.MerchantId = c.Params("merchantId")

	items, err := resource.service.Query(params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(items)
}
