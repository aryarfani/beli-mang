package merchant

import (
	"beli-mang/internal/helper"
	"beli-mang/internal/middleware"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Get("/merchants/nearby/*", middleware.Auth(), resource.queryNearby)
	app.Post("/admin/merchants", middleware.Auth(), resource.create)
	app.Get("/admin/merchants", middleware.Auth(), resource.query)
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

	merchantId, err := resource.service.Create(req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Merchant created successfully",
		"merchantId": merchantId,
	})
}

func (resource resource) query(c *fiber.Ctx) error {
	var params QueryMerchantsRequest
	_ = c.QueryParser(&params)

	merchants, err := resource.service.Query(params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(merchants)
}

func (resource resource) queryNearby(c *fiber.Ctx) error {
	var params QueryMerchantsRequest
	_ = c.QueryParser(&params)

	location := c.Params("*")
	locationSplit := strings.Split(location, ",")

	if len(locationSplit) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid location format",
		})
	}

	params.Latitude = locationSplit[0]
	params.Longitude = locationSplit[1]

	merchants, err := resource.service.QueryNearby(params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(merchants)
}
