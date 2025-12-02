package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Difficulty represents game difficulty level
// Easy: Shows airline, flight number, aircraft type, and departure airport. Domestic flights only.
// Medium: Shows airline and aircraft type only. Hides flight number and callsign. International flights within similar distance.
// Hard: Shows only aircraft position, speed, and altitude. All flight info hidden including departure airport.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// User represents a player in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Stats     UserStats          `bson:"stats" json:"stats"`
}

// UserStats tracks player statistics
type UserStats struct {
	GamesPlayed int `bson:"gamesPlayed" json:"gamesPlayed"`
	TotalScore  int `bson:"totalScore" json:"totalScore"`
	AvgScore    int `bson:"avgScore" json:"avgScore"`
	BestScore   int `bson:"bestScore" json:"bestScore"`
}

// GameSession represents a single game session
type GameSession struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID  string             `bson:"sessionId" json:"sessionId"`
	UserID     string             `bson:"userId,omitempty" json:"userId,omitempty"`
	Username   string             `bson:"username" json:"username"`
	StartedAt  time.Time          `bson:"startedAt" json:"startedAt"`
	EndedAt    *time.Time         `bson:"endedAt,omitempty" json:"endedAt,omitempty"`
	Difficulty Difficulty         `bson:"difficulty" json:"difficulty"`
	TotalScore int                `bson:"totalScore" json:"totalScore"`
	Rounds     []Round            `bson:"rounds" json:"rounds"`
	Status     string             `bson:"status" json:"status"` // "in_progress", "completed"
}

// Round represents a single round in a game
type Round struct {
	RoundNumber   int        `bson:"roundNumber" json:"roundNumber"`
	FlightID      string     `bson:"flightId" json:"flightId"`
	Flight        *Flight    `bson:"flight,omitempty" json:"flight,omitempty"`
	Departure     string     `bson:"departure" json:"departure"`
	ActualArrival string     `bson:"actualArrival" json:"actualArrival"`
	PlayerGuess   string     `bson:"playerGuess,omitempty" json:"playerGuess,omitempty"`
	PointsEarned  int        `bson:"pointsEarned" json:"pointsEarned"`
	GuessTime     float64    `bson:"guessTime" json:"guessTime"` // seconds
	Confidence    int        `bson:"confidence,omitempty" json:"confidence,omitempty"`
	StartedAt     time.Time  `bson:"startedAt" json:"startedAt"`
	CompletedAt   *time.Time `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
}

// LeaderboardEntry represents a score on the leaderboard
type LeaderboardEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Rank        int                `bson:"rank" json:"rank"`
	Username    string             `bson:"username" json:"username"`
	Difficulty  Difficulty         `bson:"difficulty" json:"difficulty"`
	TotalScore  int                `bson:"totalScore" json:"totalScore"`
	GamesPlayed int                `bson:"gamesPlayed" json:"gamesPlayed"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Flight represents flight data from Aviation Edge API
type Flight struct {
	ID            string    `json:"id"`
	ICAO24        string    `json:"icao24"`
	Callsign      string    `json:"callsign"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Altitude      float64   `json:"altitude"`      // feet
	Speed         float64   `json:"speed"`         // knots
	Direction     float64   `json:"direction"`     // degrees
	VerticalSpeed float64   `json:"verticalSpeed"` // feet per minute
	Status        string    `json:"status"`        // en-route, landed, etc.
	Departure     Airport   `json:"departure"`
	Arrival       Airport   `json:"arrival"`
	Aircraft      Aircraft  `json:"aircraft"`
	Airline       Airline   `json:"airline"`
	FlightNumber  string    `json:"flightNumber"`
	Hint          string    `json:"hint,omitempty"` // Optional hint for easy mode
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Airport represents airport information
type Airport struct {
	IATA      string  `json:"iata"`
	ICAO      string  `json:"icao"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Aircraft represents aircraft information
type Aircraft struct {
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	Model        string `json:"model"`
	Registration string `json:"registration"`
}

// Airline represents airline information
type Airline struct {
	IATA string `json:"iata"`
	ICAO string `json:"icao"`
	Name string `json:"name"`
}

// ScoreResult represents the result of scoring a guess
type ScoreResult struct {
	BasePoints      int     `json:"basePoints"`
	DifficultyMulti float64 `json:"difficultyMultiplier"`
	SpeedMulti      float64 `json:"speedMultiplier"`
	TotalPoints     int     `json:"totalPoints"`
	MatchType       string  `json:"matchType"` // exact, family, region, distance, wrong
	DistanceKm      float64 `json:"distanceKm"`
	CorrectAirport  Airport `json:"correctAirport"`
	GuessedAirport  Airport `json:"guessedAirport"`
}

// Request/Response types

// StartGameRequest represents the request to start a new game
type StartGameRequest struct {
	Username   string     `json:"username" binding:"required"`
	Difficulty Difficulty `json:"difficulty" binding:"required"`
}

// StartGameResponse represents the response when starting a game
type StartGameResponse struct {
	SessionID    string     `json:"sessionId"`
	Difficulty   Difficulty `json:"difficulty"`
	TotalRounds  int        `json:"totalRounds"`
	CurrentRound int        `json:"currentRound"`
	Flight       *Flight    `json:"flight"`
}

// GuessRequest represents a player's guess
type GuessRequest struct {
	SessionID   string `json:"sessionId" binding:"required"`
	AirportIATA string `json:"airportIata" binding:"required"`
	Confidence  int    `json:"confidence"`
}

// GuessResponse represents the response after a guess
type GuessResponse struct {
	Score       ScoreResult `json:"score"`
	RoundNumber int         `json:"roundNumber"`
	IsGameOver  bool        `json:"isGameOver"`
	NextFlight  *Flight     `json:"nextFlight,omitempty"`
	TotalScore  int         `json:"totalScore"`
}

// EndGameRequest represents the request to end a game
type EndGameRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

// EndGameResponse represents the final game results
type EndGameResponse struct {
	SessionID  string     `json:"sessionId"`
	TotalScore int        `json:"totalScore"`
	Rounds     []Round    `json:"rounds"`
	Rank       int        `json:"rank"`
	Difficulty Difficulty `json:"difficulty"`
}

// GetFlightsRequest represents query parameters for getting flights
type GetFlightsRequest struct {
	Difficulty Difficulty `form:"difficulty"`
	Limit      int        `form:"limit"`
}

// WebSocket message types

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WSFlightUpdate represents a flight position update
type WSFlightUpdate struct {
	Flights []Flight `json:"flights"`
}

// WSRoundStart represents the start of a new round
type WSRoundStart struct {
	SessionID   string  `json:"sessionId"`
	RoundNumber int     `json:"roundNumber"`
	Flight      *Flight `json:"flight"`
}

// WSGuessResult represents the result of a guess
type WSGuessResult struct {
	SessionID   string      `json:"sessionId"`
	RoundNumber int         `json:"roundNumber"`
	Score       ScoreResult `json:"score"`
	TotalScore  int         `json:"totalScore"`
}

// WSGameEnd represents the end of a game
type WSGameEnd struct {
	SessionID  string `json:"sessionId"`
	TotalScore int    `json:"totalScore"`
	Rank       int    `json:"rank"`
}
