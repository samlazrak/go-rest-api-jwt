package main

// if error change to v1 ( remove /v2 )
import (
	"github.com/gofiber/fiber/v2"
	"jwt/database"
	"jwt/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":9000")
}

