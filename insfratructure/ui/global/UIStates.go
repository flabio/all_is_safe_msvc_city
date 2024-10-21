package global

import "github.com/gofiber/fiber/v2"

type UIStates interface {
	GetStatesFindAll(c *fiber.Ctx) error
	GetStatesFindById(c *fiber.Ctx) error
	GetStatesFindByIdOfCity(c *fiber.Ctx) error
	CreateState(c *fiber.Ctx) error
	UpdateState(c *fiber.Ctx) error
	DeleteState(c *fiber.Ctx) error
}
