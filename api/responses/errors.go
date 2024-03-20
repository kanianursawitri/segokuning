package responses

import "github.com/gofiber/fiber/v2"

func ErrorBadRequest(c *fiber.Ctx, m string) error {
	return c.Status(400).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}

func ErrorConflict(c *fiber.Ctx, m string) error {
	return c.Status(409).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}

func ErrorNotFound(c *fiber.Ctx, m string) error {
	return c.Status(404).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}

func ErrorInternalServerError(c *fiber.Ctx, m string) error {
	return c.Status(500).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}

func ErrorForbidden(c *fiber.Ctx, m string) error {
	return c.Status(403).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}

func ErrorUnauthorized(c *fiber.Ctx, m string) error {
	return c.Status(401).JSON(map[string]interface{}{
		"status":  "Error",
		"message": m,
	})
}
