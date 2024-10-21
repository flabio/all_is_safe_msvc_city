package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/handler"
)

var (
	hadlerCity = handler.NewCityHandler()
)

func NewCityRouter(app *fiber.App) {
	api := app.Group("/api/cities")
	api.Get("/", func(c *fiber.Ctx) error {
		return hadlerCity.GetCityFindAll(c)
	}).Get("/:id", func(c *fiber.Ctx) error {
		return hadlerCity.GetCityFindById(c)
	}).Post("/", func(c *fiber.Ctx) error {
		return hadlerCity.CreateCity(c)
	}).Put("/:id", func(c *fiber.Ctx) error {
		return hadlerCity.UpdateCity(c)
	}).Delete("/:id", func(c *fiber.Ctx) error {
		return hadlerCity.DeleteCity(c)
	})

}
