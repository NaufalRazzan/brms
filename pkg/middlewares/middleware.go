package middlewares

import "github.com/gofiber/fiber/v2"

func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"message": e.Message,
				})
			}
		}
		return nil
	}
}

func UndefinedRoutesMiddleware(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "resource not found",
	})
}
