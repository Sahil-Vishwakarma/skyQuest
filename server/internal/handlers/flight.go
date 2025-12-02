package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/services"
)

type FlightHandler struct {
	flightService *services.FlightService
}

func NewFlightHandler(flightService *services.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

// GetFlights handles GET /api/flights
func (h *FlightHandler) GetFlights(c *gin.Context) {
	var req models.GetFlightsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 50
	}

	flights := h.flightService.GetFlights(req.Difficulty)

	// Limit results
	if len(flights) > req.Limit {
		flights = flights[:req.Limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"flights": flights,
		"count":   len(flights),
	})
}

// GetAirports handles GET /api/airports
func (h *FlightHandler) GetAirports(c *gin.Context) {
	airports := h.flightService.GetAllAirports()
	c.JSON(http.StatusOK, gin.H{
		"airports": airports,
		"count":    len(airports),
	})
}

