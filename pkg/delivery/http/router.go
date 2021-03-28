package http

import (
	"amazing_talker/configuration"

	"github.com/labstack/echo/v4"
)

// SetRoutes ...
func SetRoutes(e *echo.Echo, h *Handler, cfg *configuration.App) {
	rootV1 := e.Group("/v1")

	{
		auth := rootV1.Group("/user")
		auth.POST("/register", h.Register)
	}

}
