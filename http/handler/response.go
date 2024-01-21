package handler

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Data:    nil,
		Error:   err,
	})
}
