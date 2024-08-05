package handlers

import (
	"github.com/NikoMalik/GoTrack/logEvent"
	"github.com/NikoMalik/GoTrack/util"
	"github.com/NikoMalik/GoTrack/views/errorsTempl"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	err error
}

func (e Error) Error() string {
	return e.err.Error()
}

func AppError(err error) Error {
	return Error{
		err: err,
	}
}

// Handle404 handles 404 errors by rendering the 404 error template.
func Handle404(c *fiber.Ctx) error {
	if err := Render(c, errorsTempl.Error404(c)); err != nil {
		return err
	}
	return nil
}

// Handle500 handles 500 errors by rendering the 500 error template.
func Handle500(c *fiber.Ctx) error {
	if err := Render(c, errorsTempl.Error500(c)); err != nil {
		return err
	}
	return nil
}

// ErrorHandler is a custom error handler for Fiber.
// ErrorHandler is a custom error handler for Fiber.
func ErrorHandler(c *fiber.Ctx, err error) error {
	logEvent.Log("error", err.Error())
	if err == fiber.ErrNotFound {
		return Handle404(c)
	} else if err == fiber.ErrInternalServerError || util.IsErrNoRecords(err) {
		return Handle500(c)
	} else {
		// Pass the error message to the context and redirect back
		c.Locals("appError", err.Error())
		return c.Redirect(c.Get("Referer", "/"))
	}
}
