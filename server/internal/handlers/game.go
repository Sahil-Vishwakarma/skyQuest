package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/services"
)

type GameHandler struct {
	gameService  *services.GameService
	scoreService *services.ScoreService
}

func NewGameHandler(gameService *services.GameService, scoreService *services.ScoreService) *GameHandler {
	return &GameHandler{
		gameService:  gameService,
		scoreService: scoreService,
	}
}

// StartGame handles POST /api/game/start
func (h *GameHandler) StartGame(c *gin.Context) {
	var req models.StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate difficulty
	switch req.Difficulty {
	case models.DifficultyEasy, models.DifficultyMedium, models.DifficultyHard:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty. Must be: easy, medium, or hard"})
		return
	}

	resp, err := h.gameService.StartGame(c.Request.Context(), req)
	if err != nil {
		if err == services.ErrNoFlights {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No flights available. Please try again later."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start game: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SubmitGuess handles POST /api/game/guess
func (h *GameHandler) SubmitGuess(c *gin.Context) {
	var req models.GuessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	resp, err := h.gameService.SubmitGuess(c.Request.Context(), req)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		case services.ErrGameCompleted:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Game already completed"})
		case services.ErrInvalidRound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid round"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process guess: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// EndGame handles POST /api/game/end
func (h *GameHandler) EndGame(c *gin.Context) {
	var req models.EndGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get the session to save score
	session, err := h.gameService.GetSession(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		return
	}

	// End the game
	resp, err := h.gameService.EndGame(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end game: " + err.Error()})
		return
	}

	// Save score to leaderboard
	if err := h.scoreService.SaveScore(c.Request.Context(), session); err != nil {
		// Log but don't fail the request
		// Score saving is non-critical
	}

	// Get user's rank
	rank, err := h.scoreService.GetUserRank(c.Request.Context(), session.Username, session.Difficulty)
	if err == nil {
		resp.Rank = rank
	}

	c.JSON(http.StatusOK, resp)
}

