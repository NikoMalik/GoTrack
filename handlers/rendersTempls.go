package handlers

import (
	"os"

	"github.com/NikoMalik/GoTrack/views/layouts"
	"github.com/gofiber/fiber/v2"
)

func HandleGetHome(c *fiber.Ctx) error {
	// accessToken := c.Query("access_Token")
	// if len(accessToken) == 0 {
	// 	Render(c, layouts.CallbackScript())
	// 	log.Println("no access token")
	// 	return c.Redirect("/auth/callback/" + accessToken)
	// }

	// if len(accessToken) > 0 {

	// 	setAuthCookie(c, accessToken)

	// }

	if err := Render(c, layouts.App()); err != nil {
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

	if err := Render(c, layouts.Pricing(context)); err != nil {
		return err
	}

	return nil
}

// HandleGetSignup renders the signup page
func HandleGetSignup(c *fiber.Ctx) error {
	// Создаем пустые значения формы и ошибки

	if err := Render(c, layouts.SignupIndex()); err != nil {
		return err
	}

	return nil
}

func HandleGetLogin(c *fiber.Ctx) error {
	if err := Render(c, layouts.LoginIndex()); err != nil {
		return err
	}
	return nil
}
