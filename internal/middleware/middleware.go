package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

// Logger middleware logs all HTTP requests. It runs after otelfiber (see
// routes.SetupRoutes), so c.UserContext() carries the request's span; when
// present, its trace/span IDs are included so log lines can be correlated
// with the corresponding X-Ray trace.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		spanCtx := trace.SpanContextFromContext(c.UserContext())
		if spanCtx.IsValid() {
			log.Printf(
				"%s %s - %d - %v - trace_id=%s span_id=%s",
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				time.Since(start),
				spanCtx.TraceID(),
				spanCtx.SpanID(),
			)
		} else {
			log.Printf(
				"%s %s - %d - %v",
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				time.Since(start),
			)
		}

		return err
	}
}

// ErrorHandler handles errors and returns JSON responses
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}
