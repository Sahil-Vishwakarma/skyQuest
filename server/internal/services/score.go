package services

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreService struct {
	repo *repository.MongoRepository
	// In-memory leaderboard fallback when MongoDB is unavailable
	memoryScores    map[string]*models.LeaderboardEntry // key: username:difficulty
	memoryScoresMux sync.RWMutex
}

func NewScoreService(repo *repository.MongoRepository) *ScoreService {
	return &ScoreService{
		repo:         repo,
		memoryScores: make(map[string]*models.LeaderboardEntry),
	}
}

// SaveScore saves the final game score to leaderboard
func (s *ScoreService) SaveScore(ctx context.Context, session *models.GameSession) error {
	if s.repo != nil {
		return s.repo.SaveScore(ctx, session)
	}
	
	// In-memory fallback
	s.memoryScoresMux.Lock()
	defer s.memoryScoresMux.Unlock()
	
	key := session.Username + ":" + string(session.Difficulty)
	existing, ok := s.memoryScores[key]
	
	if !ok || session.TotalScore > existing.TotalScore {
		s.memoryScores[key] = &models.LeaderboardEntry{
			ID:          primitive.NewObjectID(),
			Username:    session.Username,
			Difficulty:  session.Difficulty,
			TotalScore:  session.TotalScore,
			GamesPlayed: 1,
			UpdatedAt:   time.Now(),
		}
		if ok {
			s.memoryScores[key].GamesPlayed = existing.GamesPlayed + 1
		}
	} else if ok {
		existing.GamesPlayed++
		existing.UpdatedAt = time.Now()
	}
	
	return nil
}

// GetLeaderboard retrieves top scores
func (s *ScoreService) GetLeaderboard(ctx context.Context, difficulty models.Difficulty, limit int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	
	if s.repo != nil {
		return s.repo.GetLeaderboard(ctx, difficulty, limit)
	}
	
	// In-memory fallback
	s.memoryScoresMux.RLock()
	defer s.memoryScoresMux.RUnlock()
	
	var entries []models.LeaderboardEntry
	for _, entry := range s.memoryScores {
		if difficulty == "" || entry.Difficulty == difficulty {
			entries = append(entries, *entry)
		}
	}
	
	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TotalScore > entries[j].TotalScore
	})
	
	// Assign ranks and limit
	for i := range entries {
		entries[i].Rank = i + 1
	}
	
	if len(entries) > limit {
		entries = entries[:limit]
	}
	
	return entries, nil
}

// GetUserRank gets a user's rank for a specific difficulty
func (s *ScoreService) GetUserRank(ctx context.Context, username string, difficulty models.Difficulty) (int, error) {
	if s.repo != nil {
		return s.repo.GetUserRank(ctx, username, difficulty)
	}
	
	// In-memory fallback
	entries, _ := s.GetLeaderboard(ctx, difficulty, 100)
	for _, entry := range entries {
		if entry.Username == username {
			return entry.Rank, nil
		}
	}
	return 0, nil
}
