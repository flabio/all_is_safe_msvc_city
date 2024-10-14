package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/handler"
)

var (
	hadlerStates = handler.NewStatesHandler()
)

func NewStatesRouter(app *fiber.App) {
	api := app.Group("/api/states")
	api.Get("/", func(c *fiber.Ctx) error {
		return hadlerStates.GetStatesFindAll(c)
	}).Get("/:id", func(c *fiber.Ctx) error {
		return hadlerStates.GetStatesFindById(c)
	}).Get("/city/:id", func(c *fiber.Ctx) error {
		return hadlerStates.GetStatesFindByIdOfCity(c)
	}).Post("/", func(c *fiber.Ctx) error {
		return hadlerStates.CreateState(c)
	}).Put("/:id", func(c *fiber.Ctx) error {
		return hadlerStates.UpdateState(c)
	}).Delete("/:id", func(c *fiber.Ctx) error {
		return hadlerStates.DeleteState(c)
	})

}
