package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/NikoMalik/GoTrack/db"
	"github.com/NikoMalik/GoTrack/goaster"
	"github.com/NikoMalik/GoTrack/handlers"
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

type Client struct {
	*websocket.Conn
}

var wsClient Client

var app = fiber.New(fiber.Config{

	JSONEncoder: json.Marshal,
	JSONDecoder: json.Unmarshal,

	// Override default error handler
	ErrorHandler: handlers.ErrorHandler,
})

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in main", r)
			debug.PrintStack()
		}
	}()

	app.Use(middleware.WithAuthUser)

	initEverything()
	app.Static("static", "./static", fiber.Static{
		Compress:      true,
		CacheDuration: 0,
	})

	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	// app.Use(middleware.WithAuthenticatedUser)

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

	log.Fatal(app.Listen(":8000"))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}
	sb.Init()
	db.Init()
	goaster.Init()

	defer db.Bun.Close()
	logEvent.Init("GO_TRACK_LOG")
	return nil
}
