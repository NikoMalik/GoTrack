package middleware

import (
	"log"
	"strings"
	"sync/atomic"

	"github.com/NikoMalik/GoTrack/data"
	"github.com/NikoMalik/GoTrack/logEvent"
	"github.com/NikoMalik/GoTrack/sb"
	"github.com/gofiber/fiber/v2"
)

var (
	authenticatedUser atomic.Value // auth userModel &
)

// WithAuthUser is a middleware that checks for a valid user in the "access_Token" cookie.
func WithAuthUser(c *fiber.Ctx) error {
	// Skip middleware for public paths
	if strings.Contains(c.Path(), "/static") {
		return c.Next()
	}
	log.Println("with auth user")

	// Get the cookie "access_Token"
	cookie := c.Cookies("access_Token")
	if cookie == "" {

		clearAuthenticatedUser()
		return c.Next()
	}

	if len(c.Cookies("access_Token")) == 0 {
		return c.Next()
	}

	// Verify user through Supabase
	resp, err := sb.Client.Auth.User(c.Context(), cookie)
	if err != nil {
		logEvent.Log("error", "authentication error", "err", "probably invalid access token")
		c.ClearCookie("access_Token")

		clearAuthenticatedUser()
		return c.Redirect("/")

	}

	name := ""
	if n, ok := resp.UserMetadata["Name"].(string); ok {
		name = n
	} else if u, ok := resp.UserMetadata["user_name"].(string); ok {
		name = u
	}

	// Initialize user information
	user := &data.AuthenticatedUser{
		LoggedIn: true,
		Email:    resp.Email,
		Name:     name,
	}

	c.Locals("user", user)
	setAuthenticatedUser(user)

	// Proceed to the next handler
	return c.Next()
}

// setAuthenticatedUser safely sets the authenticated user
func setAuthenticatedUser(user *data.AuthenticatedUser) {
	authenticatedUser.Store(user)
}

// clearAuthenticatedUser safely clears the authenticated user
func clearAuthenticatedUser() {
	authenticatedUser.Store(&data.AuthenticatedUser{})
}

// GetAuthenticatedUser safely retrieves the authenticated user
func GetAuthenticatedUser() *data.AuthenticatedUser {
	if user, ok := authenticatedUser.Load().(*data.AuthenticatedUser); ok {
		return user
	}
	return &data.AuthenticatedUser{}
}

func WithAuthenticatedUser(c *fiber.Ctx) error {
	user := &data.AuthenticatedUser{}
	c.Locals("user", user)

	if len(c.Cookies("access_Token")) == 0 {
		return c.Next()
	}

	_, err := sb.Client.Auth.User(c.Context(), c.Cookies("access_Token"))
	if err != nil {
		logEvent.Log("error", "authentication error", "err", "probably invalid access token")
		c.ClearCookie("access_Token")

		return c.Redirect("/")
	}

	ourUser := &data.AuthenticatedUser{ID: user.ID, Email: user.Email}
	logEvent.Log("msg", "user authenticated", "email", ourUser.Email)
	c.Locals("user", ourUser)

	return c.Next()
}

func IfUserAuth(c *fiber.Ctx) error {
	if GetAuthenticatedUser().LoggedIn == true {
		return c.Redirect("/")
	}
	return c.Next()
}

func IfUserNotAuth(c *fiber.Ctx) error {
	if GetAuthenticatedUser().LoggedIn == false {
		return c.Redirect("/")
	}
	return c.Next()
}
