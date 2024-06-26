package responses

import (
	"segokuning/db/entity"

	"github.com/gofiber/fiber/v2"
)

type TheResponse struct {
	StatusCode  int         `json:"code"`
	StatusError bool        `json:"error"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}

type TheResponseUpload struct {
	StatusCode  int    `json:"code"`
	StatusError bool   `json:"error"`
	Message     string `json:"message"`
	Filename    string `json:"filename"`
}

type TheResponseCount struct {
	StatusCode  int         `json:"code"`
	StatusError bool        `json:"error"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	Count       int         `json:"count"`
}

func ReturnTheResponse(c *fiber.Ctx, se bool, sc int, m string, dt interface{}) error {
	tr := TheResponse{sc, se, m, dt}

	return c.Status(sc).JSON(tr)
}

func ReturnTheResponseCount(c *fiber.Ctx, se bool, sc int, m string, dt interface{}, ct int) error {
	tr := TheResponseCount{sc, se, m, dt, ct}

	return c.Status(sc).JSON(tr)
}

func ReturnTheResponseUpload(c *fiber.Ctx, se bool, sc int, m string, f string) error {
	tr := TheResponseUpload{sc, se, m, f}

	return c.Status(sc).JSON(tr)
}

func Success(c *fiber.Ctx, m interface{}) error {
	return c.Status(200).JSON(map[string]interface{}{
		"status": "Success",
		"data":   m,
	})
}

func SuccessMeta(c *fiber.Ctx, m interface{}, meta entity.Meta) error {
	return c.Status(200).JSON(map[string]interface{}{
		"status": "Success",
		"data":   m,
		"meta":   meta,
	})
}
