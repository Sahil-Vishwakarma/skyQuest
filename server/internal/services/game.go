package services

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skyquest/server/internal/models"
	"github.com/skyquest/server/internal/repository"
	"github.com/skyquest/server/pkg/hints"
)

const (
	TotalRounds = 10
)

var (
	ErrSessionNotFound = errors.New("game session not found")
	ErrGameCompleted   = errors.New("game already completed")
	ErrInvalidRound    = errors.New("invalid round")
	ErrNoFlights       = errors.New("no flights available")
)

type GameService struct {
	repo          *repository.MongoRepository
	flightService *FlightService
	// In-memory session storage (fallback when MongoDB is unavailable)
	sessions    map[string]*models.GameSession
	sessionsMux sync.RWMutex
}

func NewGameService(repo *repository.MongoRepository, flightService *FlightService) *GameService {
	return &GameService{
		repo:          repo,
		flightService: flightService,
		sessions:      make(map[string]*models.GameSession),
	}
}

// StartGame creates a new game session
func (s *GameService) StartGame(ctx context.Context, req models.StartGameRequest) (*models.StartGameResponse, error) {
	// Get random flights for the game based on difficulty
	flights := s.flightService.GetRandomFlights(req.Difficulty, TotalRounds)
	if len(flights) == 0 {
		return nil, ErrNoFlights
	}

	// Ensure we have enough flights, duplicate if necessary
	for len(flights) < TotalRounds {
		flights = append(flights, flights[len(flights)%len(flights)])
	}

	sessionID := uuid.New().String()
	now := time.Now()

	// Create rounds
	rounds := make([]models.Round, TotalRounds)
	for i := 0; i < TotalRounds; i++ {
		flight := flights[i]
		rounds[i] = models.Round{
			RoundNumber:   i + 1,
			FlightID:      flight.ID,
			Flight:        &flight,
			Departure:     flight.Departure.IATA,
			ActualArrival: flight.Arrival.IATA,
			StartedAt:     now,
		}
	}

	session := &models.GameSession{
		SessionID:  sessionID,
		Username:   req.Username,
		StartedAt:  now,
		Difficulty: req.Difficulty,
		TotalScore: 0,
		Rounds:     rounds,
		Status:     "in_progress",
	}

	// Store session (MongoDB or in-memory fallback)
	if err := s.createSession(ctx, session); err != nil {
		return nil, err
	}

	// Prepare first flight for response (hide destination based on difficulty)
	firstFlight := s.prepareFlightForDisplay(flights[0], req.Difficulty)

	return &models.StartGameResponse{
		SessionID:    sessionID,
		Difficulty:   req.Difficulty,
		TotalRounds:  TotalRounds,
		CurrentRound: 1,
		Flight:       &firstFlight,
	}, nil
}

// SubmitGuess processes a player's guess
func (s *GameService) SubmitGuess(ctx context.Context, req models.GuessRequest) (*models.GuessResponse, error) {
	session, err := s.getSession(ctx, req.SessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.Status == "completed" {
		return nil, ErrGameCompleted
	}

	// Find current round (first incomplete round)
	var currentRound *models.Round
	var roundIndex int
	for i := range session.Rounds {
		if session.Rounds[i].PlayerGuess == "" {
			currentRound = &session.Rounds[i]
			roundIndex = i
			break
		}
	}

	if currentRound == nil {
		return nil, ErrInvalidRound
	}

	// Calculate score
	guessTime := time.Since(currentRound.StartedAt).Seconds()

	// Get actual airport from the stored flight data
	var actualAirport *models.Airport
	if currentRound.Flight != nil {
		actualAirport = &currentRound.Flight.Arrival
	}

	score := s.calculateScore(
		currentRound.ActualArrival,
		req.AirportIATA,
		session.Difficulty,
		guessTime,
		actualAirport,
	)

	// Update round
	now := time.Now()
	session.Rounds[roundIndex].PlayerGuess = req.AirportIATA
	session.Rounds[roundIndex].PointsEarned = score.TotalPoints
	session.Rounds[roundIndex].GuessTime = guessTime
	session.Rounds[roundIndex].Confidence = req.Confidence
	session.Rounds[roundIndex].CompletedAt = &now

	// Update total score
	session.TotalScore += score.TotalPoints

	// Check if game is over
	isGameOver := roundIndex == TotalRounds-1

	if isGameOver {
		session.Status = "completed"
		session.EndedAt = &now
	}

	// Update session
	if err := s.updateSession(ctx, session); err != nil {
		return nil, err
	}

	// Prepare next flight if game continues
	var nextFlight *models.Flight
	if !isGameOver && roundIndex+1 < len(session.Rounds) {
		nextRound := session.Rounds[roundIndex+1]
		nextRound.StartedAt = now
		session.Rounds[roundIndex+1] = nextRound
		s.updateSession(ctx, session)

		if nextRound.Flight != nil {
			prepared := s.prepareFlightForDisplay(*nextRound.Flight, session.Difficulty)
			nextFlight = &prepared
		}
	}

	return &models.GuessResponse{
		Score:       score,
		RoundNumber: currentRound.RoundNumber,
		IsGameOver:  isGameOver,
		NextFlight:  nextFlight,
		TotalScore:  session.TotalScore,
	}, nil
}

// EndGame finalizes a game session
func (s *GameService) EndGame(ctx context.Context, sessionID string) (*models.EndGameResponse, error) {
	session, err := s.getSession(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Mark as completed if not already
	if session.Status != "completed" {
		now := time.Now()
		session.Status = "completed"
		session.EndedAt = &now
		if err := s.updateSession(ctx, session); err != nil {
			return nil, err
		}
	}

	// Get rank (will be calculated after saving score)
	rank := 0

	return &models.EndGameResponse{
		SessionID:  session.SessionID,
		TotalScore: session.TotalScore,
		Rounds:     session.Rounds,
		Rank:       rank,
		Difficulty: session.Difficulty,
	}, nil
}

// GetSession retrieves a game session
func (s *GameService) GetSession(ctx context.Context, sessionID string) (*models.GameSession, error) {
	return s.getSession(ctx, sessionID)
}

// Helper methods for session storage with fallback

func (s *GameService) createSession(ctx context.Context, session *models.GameSession) error {
	if s.repo != nil {
		return s.repo.CreateSession(ctx, session)
	}
	// In-memory fallback
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()
	s.sessions[session.SessionID] = session
	return nil
}

func (s *GameService) getSession(ctx context.Context, sessionID string) (*models.GameSession, error) {
	if s.repo != nil {
		return s.repo.GetSession(ctx, sessionID)
	}
	// In-memory fallback
	s.sessionsMux.RLock()
	defer s.sessionsMux.RUnlock()
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *GameService) updateSession(ctx context.Context, session *models.GameSession) error {
	if s.repo != nil {
		return s.repo.UpdateSession(ctx, session)
	}
	// In-memory fallback
	s.sessionsMux.Lock()
	defer s.sessionsMux.Unlock()
	s.sessions[session.SessionID] = session
	return nil
}

// calculateScore determines points based on guess accuracy
func (s *GameService) calculateScore(actualIATA, guessedIATA string, difficulty models.Difficulty, guessTime float64, actualAirportInfo *models.Airport) models.ScoreResult {
	result := models.ScoreResult{
		DifficultyMulti: getDifficultyMultiplier(difficulty),
		SpeedMulti:      getSpeedMultiplier(guessTime),
	}

	// Use provided actual airport info, or look up from database
	var actualAirport models.Airport
	var actualOK bool
	if actualAirportInfo != nil && actualAirportInfo.IATA != "" && actualAirportInfo.IATA != "???" {
		actualAirport = *actualAirportInfo
		actualOK = true
	} else {
		actualAirport, actualOK = s.flightService.GetAirport(actualIATA)
	}

	guessedAirport, guessedOK := s.flightService.GetAirport(guessedIATA)

	if actualOK {
		result.CorrectAirport = actualAirport
	}
	if guessedOK {
		result.GuessedAirport = guessedAirport
	}

	// Exact match
	if actualIATA == guessedIATA {
		result.BasePoints = 1000
		result.MatchType = "exact"
		result.DistanceKm = 0
	} else if actualOK && guessedOK {
		// Calculate distance
		distance := CalculateDistance(
			actualAirport.Latitude, actualAirport.Longitude,
			guessedAirport.Latitude, guessedAirport.Longitude,
		)
		result.DistanceKm = distance

		// Airport family (same city)
		if actualAirport.City == guessedAirport.City {
			result.BasePoints = 750
			result.MatchType = "family"
		} else if actualAirport.Country == guessedAirport.Country {
			// Same country
			result.BasePoints = 500
			result.MatchType = "country"
		} else if distance <= 500 {
			// Within 500km
			result.BasePoints = 250
			result.MatchType = "distance"
		} else {
			// Wrong
			result.BasePoints = 0
			result.MatchType = "wrong"
		}
	} else {
		result.BasePoints = 0
		result.MatchType = "wrong"
	}

	// Calculate total with multipliers
	result.TotalPoints = int(float64(result.BasePoints) * result.DifficultyMulti * result.SpeedMulti)

	return result
}

func getDifficultyMultiplier(difficulty models.Difficulty) float64 {
	switch difficulty {
	case models.DifficultyEasy:
		return 1.0
	case models.DifficultyMedium:
		return 1.5
	case models.DifficultyHard:
		return 2.0
	default:
		return 1.0
	}
}

func getSpeedMultiplier(guessTime float64) float64 {
	if guessTime <= 10 {
		return 1.3
	} else if guessTime <= 30 {
		return 1.1
	}
	return 1.0
}

// prepareFlightForDisplay hides certain info based on difficulty and generates random flight position
func (s *GameService) prepareFlightForDisplay(flight models.Flight, difficulty models.Difficulty) models.Flight {
	displayFlight := flight

	// Store airport info before hiding
	arrivalCity := flight.Arrival.City
	arrivalLat := flight.Arrival.Latitude
	arrivalLon := flight.Arrival.Longitude
	departureLat := flight.Departure.Latitude
	departureLon := flight.Departure.Longitude

	// Default to a central location if departure coordinates are missing
	if departureLat == 0 && departureLon == 0 {
		departureLat = 40.0  // Default latitude
		departureLon = -74.0 // Default longitude (around NYC area)
	}

	// Use departure for arrival fallback if needed (for heading calculation)
	if arrivalLat == 0 && arrivalLon == 0 {
		arrivalLat = departureLat + 5.0 // Offset so heading isn't 0
		arrivalLon = departureLon + 5.0
	}

	// ALWAYS generate random aircraft position CENTERED around the DEPARTURE/ORIGIN airport
	// Generate small random offset (0.3-1.0 degrees away, ~33-111 km) for close proximity
	distance := 0.3 + rand.Float64()*0.7  // 0.3 to 1.0 degrees (~33-111 km)
	angle := rand.Float64() * 2 * math.Pi // Random angle in radians

	// Calculate aircraft position at random offset from DEPARTURE airport (centered around it)
	displayFlight.Latitude = departureLat + distance*math.Cos(angle)
	displayFlight.Longitude = departureLon + distance*math.Sin(angle)

	// Calculate heading from aircraft position towards the ARRIVAL airport
	displayFlight.Direction = calculateBearing(
		displayFlight.Latitude, displayFlight.Longitude,
		arrivalLat, arrivalLon,
	)

	// Randomize altitude and speed within realistic ranges
	displayFlight.Altitude = 28000 + rand.Float64()*10000    // 28,000-38,000 ft
	displayFlight.Speed = 420 + rand.Float64()*100           // 420-520 knots
	displayFlight.VerticalSpeed = -500 + rand.Float64()*1000 // -500 to +500 ft/min

	// Always hide the actual arrival for guessing
	displayFlight.Arrival = models.Airport{
		IATA: "???",
		ICAO: "????",
		Name: "Unknown Destination",
	}

	switch difficulty {
	case models.DifficultyEasy:
		// Show most info - flight number visible
		// Add a hint fact about the destination city
		if arrivalCity != "" {
			displayFlight.Hint = hints.GetCityFact(arrivalCity, flight.ID)
		}
	case models.DifficultyMedium:
		// Hide flight number, show airline and aircraft
		displayFlight.FlightNumber = ""
		displayFlight.Callsign = ""
	case models.DifficultyHard:
		// Hide almost everything except position and speed
		displayFlight.FlightNumber = ""
		displayFlight.Callsign = ""
		displayFlight.Airline = models.Airline{}
		displayFlight.Departure = models.Airport{
			IATA: "???",
			ICAO: "????",
			Name: "Unknown Origin",
		}
	}

	return displayFlight
}

// calculateBearing calculates the bearing/heading from point 1 to point 2 in degrees
func calculateBearing(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Calculate bearing
	x := math.Sin(deltaLon) * math.Cos(lat2Rad)
	y := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(deltaLon)

	bearing := math.Atan2(x, y)

	// Convert to degrees and normalize to 0-360
	bearingDeg := bearing * 180 / math.Pi
	if bearingDeg < 0 {
		bearingDeg += 360
	}

	return bearingDeg
}
