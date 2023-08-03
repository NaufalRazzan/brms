package middlewares

import "github.com/gofiber/fiber/v2"

func ErrorMiddleware(c *fiber.Ctx, err error) error{
	if e, ok := err.(*fiber.Error); ok{
		return c.Status(e.Code).JSON(fiber.Map{
			"message": e.Message,
		})
	}

	return c.SendStatus(fiber.StatusInternalServerError)
}

func UndefinedRoutesMiddleware(c *fiber.Ctx) error{
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "resource not found",
	})
}