package authRouter

import (
	"github.com/NikoMalik/GoTrack/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/signup", handlers.HandleGetSignup) // get page for signup
	auth.Post("/signup", handlers.HandleSignupWithEmail)
	auth.Get("/login", handlers.HandleGetLogin)
	auth.Post("/login", handlers.HandleLoginWithEmail)
	auth.Post("/resend-email-verification", handlers.HandleResendVerificationCode)
	// auth.Post("/signup/google", handlers.HandleSignInWithGoogle)
	auth.Post("/signup/github", handlers.HandleSignInWithGithub)
	auth.Get("/callback/", handlers.HandleAuthCallback)
	auth.Post("/callback", handlers.HandleAuthCallback)
	auth.Get("/signout", handlers.HandleGetSignOut) // just leave with token
	auth.Get("/signin", handlers.HandleGetSignIn)   //get page

}
