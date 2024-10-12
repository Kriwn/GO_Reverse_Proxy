package router

import (
	"github.com/gofiber/fiber/v2"
)

func GetCloudfare(c *fiber.Ctx) string{
    return c.Params("*")
}

