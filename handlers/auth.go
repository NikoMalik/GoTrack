package handlers

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/NikoMalik/GoTrack/data"
	"github.com/NikoMalik/GoTrack/db"
	"github.com/NikoMalik/GoTrack/event"
	"github.com/NikoMalik/GoTrack/logEvent"
	"github.com/NikoMalik/GoTrack/sb"
	v "github.com/NikoMalik/GoTrack/validate"
	"github.com/golang-jwt/jwt/v5"

	"github.com/NikoMalik/GoTrack/views/layouts"
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
)

var signupSchema = v.Schema{
	"email": v.Rules(v.Email),
	"password": v.Rules(
		v.ContainsSpecial,
		v.ContainsUpper,
		v.Min(7),
		v.Max(50),
	),
	"name": v.Rules(v.Min(2), v.Max(50)),
}

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleSignupWithEmail(c *fiber.Ctx) error {
	params := layouts.SignupParams{
		Email:                c.FormValue("email"),
		Name:                 c.FormValue("name"),
		Password:             c.FormValue("password"),
		PasswordConfirmation: c.FormValue("passwordConfirmation"),
	}

	// Debug: Log received parameters
	log.Printf("Received params: %+v", &params)

	// Validate parameters
	errors, ok := v.Request(c.Context(), &params, signupSchema)
	if !ok {
		log.Printf("Validation errors: %+v", errors)
		return Render(c, layouts.SignupForm(params, errors))
	}

	// Sign up the user
	resp, err := sb.Client.Auth.SignUp(c.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
		Data: map[string]string{
			"Name": params.Name,
		},
	})

	if err != nil {

		log.Printf("SignUp error: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Signup failed")
	}

	// Debug: Log the response from Supabase
	log.Printf("Supabase response: %+v", resp)

	event.Emit(data.UserSignupEvent, data.UserWithSignup{
		User: resp,
	})

	logEvent.Log("msg", "user signup with email", "id", resp.ID)

	// Email sent automatically by Supabase; just render the success page
	return Render(c, layouts.SignupSuccess(resp))
}

func HandleLoginWithEmail(c *fiber.Ctx) error {
	params := supabase.UserCredentials{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	log.Printf("Received params: %+v", &params)

	errors, ok := v.Request(c.Context(), &params, authSchema)

	if !ok {
		log.Printf("Validation errors: %+v", errors)
		return Render(c, layouts.LoginForm(params, errors))
	}

	resp, err := sb.Client.Auth.SignIn(c.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})

	if err != nil {
		log.Printf("Login error: %v", err)
		if err.Error() == "invalid_grant: Email not confirmed" {
			return Render(c, layouts.Toast("Login Error", "Please confirm your email address before logging in."))
		}

		return Render(c, layouts.Toast("Login Error", "Please check your credentials and try again."))
	}

	if c.Query("access_Token") == "" {

		setAuthCookie(c, resp.AccessToken)
	}

	event.Emit(data.UserSignupEvent, data.UserWithVerificationToken{
		User: resp,
	})

	return HXRedirect(c, "/")
}

func setAuthCookie(c *fiber.Ctx, accessToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_Token",
		Value:    accessToken,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
}

func HandleResendVerificationCode(c *fiber.Ctx) error {
	id := c.Params("ID")

	var user data.User

	if err := db.Bun.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background()); err != nil {
		return err
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return logEvent.Log("msg", "user already verified", "id", user.ID)

	}

	token, err := createVerificationToken(c, id)

	if err != nil {
		return err
	}

	logEvent.Log("msg", "user verification email sent", "email", user.Email)

	_ = token

	return c.SendStatus(200)

}

func createVerificationToken(c *fiber.Ctx, ID string) (string, error) {
	expirystr := os.Getenv("AUTH_EMAIL_VERIFICATION_EXPIRY_IN_HOURS")
	expiry, err := strconv.Atoi(expirystr)
	if err != nil {
		expiry = 1
	}

	claims := jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return "", c.SendStatus(fiber.StatusInternalServerError)
	}

	return t, nil

}

func HandleSignInWithGoogle(c *fiber.Ctx) error { // not available yet need google developer keys

	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider: "google",
		FlowType: supabase.PKCE,
	})

	if err != nil {
		return err
	}

	return c.Redirect(resp.URL)
}

func HandleSignInWithGithub(c *fiber.Ctx) error {

	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "github",
		RedirectTo: "http://localhost:8000/auth/callback",
	})

	if err != nil {
		logEvent.Log("error", err.Error())
		return err
	}

	return HXRedirect(c, resp.URL)

}

func HandleGetSignOut(c *fiber.Ctx) error {

	if err := sb.Client.Auth.SignOut(c.Context(), c.Cookies("accessToken")); err != nil {
		return err
	}

	c.ClearCookie("access_Token")

	return c.Redirect("/")

}

// base auth func
func HandleAuthCallback(c *fiber.Ctx) error {
	accessToken := c.Query("access_Token")

	setAuthCookie(c, accessToken)

	log.Println("Access token successfully set in cookie")

	return HXRedirect(c, "/")
}

func HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("ID")

	err := sb.Client.DB.From("users").Delete().Eq("id", id).Execute(c.Context())

	if err != nil {
		c.SendStatus(500)
		return err
	}
	return c.SendStatus(200)

}
