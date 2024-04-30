package handlers

import (
	"github.com/gofiber/fiber/v3"
)

type ReadyCheckHandler struct {
}

func NewReadyCheckHandler() *ReadyCheckHandler {
	return &ReadyCheckHandler{}
}

// HealthCheck godoc
//
//	@Summary		Ready check service
//	@Description	Ready check service
//	@Tags			Health Controller
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/readyz [get]
func (hch *ReadyCheckHandler) Handle(ctx fiber.Ctx) error {
	// Add own checks
	return nil
}
