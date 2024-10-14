package global

import "github.com/gofiber/fiber/v2"

type UICity interface {
	GetCityFindAll(c *fiber.Ctx) error
	GetCityFindById(c *fiber.Ctx) error
	CreateCity(c *fiber.Ctx) error
	UpdateCity(c *fiber.Ctx) error
	DeleteCity(c *fiber.Ctx) error
}
