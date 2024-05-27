package router

import (
	"services/internal/user/controllers"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// # You can add routers in this function.
func Initialize(app *fiber.App, controller *controllers.UserController) (*fiber.App, error) {

	app.Use("/monitor", monitor.New(monitor.Config{Title: "User Service Monitoring."})) //Fiber Monitoring.

	app.Post("/register", controller.RegisterUserWithLocalParameters)

	app.Post("/login", controller.LoginUserWithLocalParameters)

	app.Get("/me", controller.GetUserWithToken)

	app.Get("/refresh-token", controller.GetTokensWithRefreshToken)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(controller.GetWebSocketHandler))

	routes := app.GetRoutes()

	app.Server().Logger.Printf("All Routes:", routes)

	return app, nil
}
