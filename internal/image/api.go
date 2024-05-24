package image

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(app *fiber.App, service Service) {
	resource := resource{service: service}

	app.Post("/image", resource.upload)
}

type resource struct {
	service Service
}

func (resource resource) upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	// Validate file
	err = resource.service.ValidateFile(fileHeader)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Upload file
	imageUrl, err := resource.service.Upload(fileHeader)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Image upload successfully",
		"data": fiber.Map{
			"imageUrl": imageUrl,
		},
	})
}
