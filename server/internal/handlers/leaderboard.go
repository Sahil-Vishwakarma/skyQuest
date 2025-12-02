package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/services"
)

type LeaderboardHandler struct {
	scoreService *services.ScoreService
}

func NewLeaderboardHandler(scoreService *services.ScoreService) *LeaderboardHandler {
	return &LeaderboardHandler{scoreService: scoreService}
}

// GetLeaderboard handles GET /api/leaderboard
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	// Get query parameters
	difficultyStr := c.DefaultQuery("difficulty", "")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var difficulty models.Difficulty
	if difficultyStr != "" {
		difficulty = models.Difficulty(difficultyStr)
		// Validate difficulty
		switch difficulty {
		case models.DifficultyEasy, models.DifficultyMedium, models.DifficultyHard:
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid difficulty. Must be: easy, medium, or hard"})
			return
		}
	}

	entries, err := h.scoreService.GetLeaderboard(c.Request.Context(), difficulty, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": entries,
		"count":       len(entries),
		"difficulty":  difficulty,
	})
}

