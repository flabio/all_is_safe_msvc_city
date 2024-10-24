package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/insfratructure/routers"
)

func main() {
	app := fiber.New()
	routers.NewCityRouter(app)
	routers.NewStatesRouter(app)
	app.Listen(":3014")
}
