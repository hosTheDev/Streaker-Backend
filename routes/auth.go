package routes

import (
	"streaker-backend/handlers"
	"streaker-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("api/auth/signup", handlers.Signup)
	app.Post("api/auth/login", handlers.Login)

	protected := app.Group("/user")
	protected.Use(middleware.JWTMiddleware())
	protected.Get("/profile", handlers.Profile)
}
