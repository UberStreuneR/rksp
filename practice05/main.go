package main

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Rofl")
	})
	app.Post("/upload", UploadFile)
	app.Static("/static", "./static")
	app.Listen(":8000")
}

func UploadFile(c *fiber.Ctx) error {
	// ctx := context.Background()
	// bucketName := os.Getenv("MINIO_BUCKET")
	file, err := c.FormFile("fileUpload")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer buffer.Close()

	dst, err := os.Create("./static/" + file.Filename)
	defer dst.Close()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if _, err := io.Copy(dst, buffer); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		// "info":  info,
		"info": "all gut",
	})
}
