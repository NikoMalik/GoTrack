package middleware

import (
	"log"
	"strings"
	"sync"

	"github.com/NikoMalik/GoTrack/data"
	"github.com/NikoMalik/GoTrack/sb"
	"github.com/gofiber/fiber/v2"
)

var (
	authenticatedUser *data.AuthenticatedUser // auth userModel &
	mu                sync.RWMutex            // mutex for escape data race
)

// WithAuthUser is a middleware that checks for a valid user in the "access_Token" cookie.
func WithAuthUser(c *fiber.Ctx) error {
	// Skip middleware for public paths
	if strings.Contains(c.Path(), "/static") {
		return c.Next()
	}

	// Get the cookie "access_Token"
	cookie := c.Cookies("access_Token")
	if cookie == "" {
		clearAuthenticatedUser()
		return c.Next()
	}

	// Verify user through Supabase
	resp, err := sb.Client.Auth.User(c.Context(), cookie)
	if err != nil {
		clearAuthenticatedUser()
		return c.Next()
	}

	// Initialize user information
	user := &data.AuthenticatedUser{
		LoggedIn: true,
	}

	// Safe extraction of Name from UserMetadata
	if name, ok := resp.UserMetadata["Name"].(string); ok {
		user.Name = name
	} else {
		// Handle the case where Name is not a string or is missing
		user.Name = "Unknown"
	}

	// Extract Email safely
	user.Email = resp.Email

	log.Println("User authenticated:", user.Email)

	// Store user in context and package-level variable
	c.Locals("user", user)
	setAuthenticatedUser(user)

	// Proceed to the next handler
	return c.Next()
}

// setAuthenticatedUser safely sets the authenticated user
func setAuthenticatedUser(user *data.AuthenticatedUser) {
	mu.Lock()
	defer mu.Unlock()
	authenticatedUser = user
}

// clearAuthenticatedUser safely clears the authenticated user
func clearAuthenticatedUser() {
	mu.Lock()
	defer mu.Unlock()
	authenticatedUser = &data.AuthenticatedUser{}
}

// GetAuthenticatedUser safely retrieves the authenticated user
func GetAuthenticatedUser() *data.AuthenticatedUser {
	mu.RLock()
	defer mu.RUnlock()
	return authenticatedUser
}
