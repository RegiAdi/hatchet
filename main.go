package main

import (
	"log"

	"github.com/RegiAdi/hatchet/bootstrap"
	"github.com/RegiAdi/hatchet/config"
	"github.com/RegiAdi/hatchet/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	bootstrap.Run()
	routes.API(app)

	err := app.Listen(":" + config.GetAppPort())

	if err != nil {
		log.Fatal("App failed to start.")
		panic(err)
	}
}
