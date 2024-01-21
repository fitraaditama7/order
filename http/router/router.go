package router

import (
	"github.com/gofiber/fiber/v2"
	"order-backend/http/handler"
	"order-backend/http/middleware"
)

type Opts struct {
	OrderHandler *handler.OrderHandler
}

func SetupRoutes(app *fiber.App, opts Opts) {
	app.Use(middleware.Logger())
	app.Get("/order", opts.OrderHandler.GetList)
}
