package handlers

import (
	"os"

	"github.com/NikoMalik/GoTrack/views/layouts"
	"github.com/gofiber/fiber/v2"
)

func HandleGetHome(c *fiber.Ctx) error {

	if cookie := c.Cookies("visited"); cookie == "" {

		c.Cookie(&fiber.Cookie{
			Name:  "visited",
			Value: "true",
		})

	}

	if err := Render(c, layouts.Index(c)); err != nil {
		return err
	}

	return nil

}

func HandleGetDashboard(c *fiber.Ctx) error {

	if err := Render(c, layouts.Dashboard(c, data)); err != nil {
		return err
	}

	return nil
}

func HandlePricing(c *fiber.Ctx) error {
	context := map[string]interface{}{
		"planFreePID":       os.Getenv("STRIPE_FREE_PID"),
		"planBusinessPID":   os.Getenv("STRIPE_BUSINESS_PID"),
		"planEnterprisePID": os.Getenv("STRIPE_ENTERPRISE_PID"),
		"starterDomains":    2,
		"businessDomains":   50,
		"enterpriseDomains": 500,
	}

	if err := Render(c, layouts.Pricing(context, c)); err != nil {
		return err
	}

	return nil
}

// HandleGetSignup renders the signup page
func HandleGetSignup(c *fiber.Ctx) error {
	// Создаем пустые значения формы и ошибки

	if err := Render(c, layouts.SignupIndex(c)); err != nil {
		return err
	}

	return nil
}

func HandleGetLogin(c *fiber.Ctx) error {
	if err := Render(c, layouts.LoginIndex(c)); err != nil {
		return err
	}
	return nil
}
