package response

import "github.com/gofiber/fiber/v2"

type UserResponse struct {
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

func StatusOK(msg string, data *fiber.Map) *UserResponse {
	return &UserResponse{
		Message: msg,
		Data:    data,
	}
}

func StatusCreated(msg string, data *fiber.Map) *UserResponse {
	return &UserResponse{
		Message: msg,
		Data:    data,
	}
}
