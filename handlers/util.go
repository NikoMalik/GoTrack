package handlers

import (
	"github.com/NikoMalik/GoTrack/data"
	"github.com/gofiber/fiber/v2"
)

func isUserSignedIn(c *fiber.Ctx) bool {
	user := getAuthenticatedUser(c)
	return user != nil
}

func HXRedirect(c *fiber.Ctx, to string) error {
	if c.Get("HX-Request") != "" {
		c.Set("HX-Redirect", to)
		c.Status(fiber.StatusFound)
		return nil
	}
	return c.Redirect(to, fiber.StatusFound)
}

func getAuthenticatedUser(c *fiber.Ctx) *data.AuthenticatedUser {
	value := c.Locals("user")
	if user, ok := value.(*data.AuthenticatedUser); ok {
		return user
	}
	return nil
}
