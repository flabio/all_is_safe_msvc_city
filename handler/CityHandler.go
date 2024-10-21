package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/insfratructure/ui/global"
	"github.com/safe_msvc_city/usecase/service"
)

type cityHandler struct {
	city global.UICity
}

func NewCityHandler() global.UICity {
	return &cityHandler{
		city: service.NewCityService(),
	}
}

func (h *cityHandler) GetCityFindAll(c *fiber.Ctx) error {
	return h.city.GetCityFindAll(c)
}
func (h *cityHandler) GetCityFindById(c *fiber.Ctx) error {
	return h.city.GetCityFindById(c)
}

func (h *cityHandler) CreateCity(c *fiber.Ctx) error {
	return h.city.CreateCity(c)
}

func (h *cityHandler) UpdateCity(c *fiber.Ctx) error {
	return h.city.UpdateCity(c)
}

func (h *cityHandler) DeleteCity(c *fiber.Ctx) error {
	return h.city.DeleteCity(c)
}
