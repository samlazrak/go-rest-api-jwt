package main

// if error change to v1 ( remove /v2 )
import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"jwt/database"
	"jwt/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	// cors allows for different port to access requests : ie app runs on 9001 vs server on 9000
	app.Use(cors.New(cors.Config{
		// With this frontend can get cookie and send it back
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":9000")
}

