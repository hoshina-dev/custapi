package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health godoc
//
//	@Summary	Liveness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})

}

// Ready godoc
//
//	@Summary	Readiness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Router		/ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	if h.db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}

	return c.JSON(fiber.Map{"statue": "ready"})
}
