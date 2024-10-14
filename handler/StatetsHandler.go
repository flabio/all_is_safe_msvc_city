package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/insfratructure/ui/global"
	"github.com/safe_msvc_city/usecase/service"
)

type statesHandler struct {
	state global.UIStates
}

func NewStatesHandler() global.UIStates {
	return &statesHandler{state: service.NewSatatesService()}
}

func (h *statesHandler) GetStatesFindAll(c *fiber.Ctx) error {
	return h.state.GetStatesFindAll(c)
}

func (h *statesHandler) GetStatesFindById(c *fiber.Ctx) error {
	return h.state.GetStatesFindById(c)
}
func (h *statesHandler) GetStatesFindByIdOfCity(c *fiber.Ctx) error {
	return h.state.GetStatesFindByIdOfCity(c)
}

func (h *statesHandler) CreateState(c *fiber.Ctx) error {
	return h.state.CreateState(c)
}

func (h *statesHandler) UpdateState(c *fiber.Ctx) error {
	return h.state.UpdateState(c)
}

func (h *statesHandler) DeleteState(c *fiber.Ctx) error {
	return h.state.DeleteState(c)
}
