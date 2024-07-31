package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/NikoMalik/GoTrack/db"
	"github.com/NikoMalik/GoTrack/logEvent"
	"github.com/NikoMalik/GoTrack/middleware"
	"github.com/NikoMalik/GoTrack/sb"

	"github.com/NikoMalik/GoTrack/router"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// var store = session.New()
// var preferenceMap map[string]string

type Client struct {
	*websocket.Conn
}

var wsClient Client

// const goTrackVersion = "1.0.0"
// const maxWorkerPoolSize = 5
// const maxJobMaxWorkers = 5

var app = fiber.New(fiber.Config{

	JSONEncoder: json.Marshal,
	JSONDecoder: json.Unmarshal,

	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError

		// Retrieve the custom status code if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		// Send custom error page
		err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
		if err != nil {
			// In case the SendFile fails
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		// Return from handler
		return nil
	},
})

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in main", r)
			debug.PrintStack()
		}
	}()
	initEverything()
	app.Static("static", "./static", fiber.Static{
		Compress:      true,
		CacheDuration: 0,
	})

	app.Use(logger.New())

	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{
	// 		JWTAlg: jwtware.RS256,
	// 		Key:    []byte(os.Getenv("JWT_SECRET_KEY")),
	// 	},
	// }))
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	app.Use(middleware.WithAuthUser)

	// app.Use(func(c *fiber.Ctx) error {
	// 	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	// 	c.Set("Pragma", "no-cache")
	// 	c.Set("Expires", "0")
	// 	c.Set("Surrogate-Control", "no-store")
	// 	return c.Next()
	// })
	// Set up routes
	router.Setup(app)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
	}()

	// Start the server
	log.Fatal(app.Listen(":8000"))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}
	sb.Init()
	db.Init()
	defer db.Bun.Close()
	logEvent.Init("GO_TRACK_LOG")
	return nil
}
