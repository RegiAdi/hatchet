package routes

import (
	"github.com/RegiAdi/hatchet/controllers"
	"github.com/RegiAdi/hatchet/middleware/shrine"
	"github.com/gofiber/fiber/v2"
)

func API(app *fiber.App) {
	api := app.Group("/api")

	api.Group("/auth")
	auth := api.Group("/auth")

	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)

	app.Use(shrine.New())

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	auth.Get("/logout", controllers.Logout)
	api.Get("/userinfo", controllers.GetUserInfo)

	api.Get("/products", controllers.GetProducts)
	api.Get("/products/:id", controllers.GetProduct)
	api.Post("/products", controllers.CreateProduct)
	api.Put("/products/:id", controllers.UpdateProduct)
	api.Delete("/products/:id", controllers.DeleteProduct)
}
