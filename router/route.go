package router

import (
	"log"

	"github.com/NikoMalik/GoTrack/handlers"
	"github.com/NikoMalik/GoTrack/routes/authRouter"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Set up authentication routes

	app.Get("/", handlers.HandleGetHome)
	app.Get("/pricing", handlers.HandlePricing)

	//auth routes

	authRouter.SetupAuthRoutes(app)

	// Set up WebSocket routes
	setupWebSocketRoutes(app)

	// errors

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return handlers.ErrorHandler(c, err)
		}
		return nil
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return handlers.Handle404(c)
	})
}

func setupWebSocketRoutes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
}
