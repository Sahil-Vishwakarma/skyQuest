export type Difficulty = 'easy' | 'medium' | 'hard'
export type GameStatus = 'idle' | 'playing' | 'finished'

export interface Airport {
  iata: string
  icao: string
  name: string
  city: string
  country: string
  latitude: number
  longitude: number
}

export interface Aircraft {
  iata: string
  icao: string
  model: string
  registration: string
}

export interface Airline {
  iata: string
  icao: string
  name: string
}

export interface Flight {
  id: string
  icao24: string
  callsign: string
  latitude: number
  longitude: number
  altitude: number
  speed: number
  direction: number
  verticalSpeed: number
  status: string
  departure: Airport
  arrival: Airport
  aircraft: Aircraft
  airline: Airline
  flightNumber: string
  hint?: string // Optional hint for easy mode
  updatedAt: string
}

export interface Round {
  roundNumber: number
  flightId: string
  flight?: Flight
  departure: string
  actualArrival: string
  playerGuess?: string
  pointsEarned: number
  guessTime: number
  confidence?: number
  startedAt: string
  completedAt?: string
}

export interface ScoreResult {
  basePoints: number
  difficultyMultiplier: number
  speedMultiplier: number
  totalPoints: number
  matchType: 'exact' | 'family' | 'country' | 'distance' | 'wrong'
  distanceKm: number
  correctAirport: Airport
  guessedAirport: Airport
}

export interface GameSession {
  id: string
  sessionId: string
  username: string
  startedAt: string
  endedAt?: string
  difficulty: Difficulty
  totalScore: number
  rounds: Round[]
  status: 'in_progress' | 'completed'
}

export interface LeaderboardEntry {
  id: string
  rank: number
  username: string
  difficulty: Difficulty
  totalScore: number
  gamesPlayed: number
  updatedAt: string
}

// API Request/Response types
export interface StartGameRequest {
  username: string
  difficulty: Difficulty
}

export interface StartGameResponse {
  sessionId: string
  difficulty: Difficulty
  totalRounds: number
  currentRound: number
  flight: Flight
}

export interface GuessRequest {
  sessionId: string
  airportIata: string
  confidence?: number
}

export interface GuessResponse {
  score: ScoreResult
  roundNumber: number
  isGameOver: boolean
  nextFlight?: Flight
  totalScore: number
}

export interface EndGameResponse {
  sessionId: string
  totalScore: number
  rounds: Round[]
  rank: number
  difficulty: Difficulty
}

export interface FlightsResponse {
  flights: Flight[]
  count: number
}

export interface LeaderboardResponse {
  leaderboard: LeaderboardEntry[]
  count: number
  difficulty: Difficulty
}

// WebSocket message types
export interface WSMessage {
  type: string
  payload: unknown
}

export interface WSFlightUpdate {
  flights: Flight[]
}

export interface WSRoundStart {
  sessionId: string
  roundNumber: number
  flight: Flight
}

export interface WSGuessResult {
  sessionId: string
  roundNumber: number
  score: ScoreResult
  totalScore: number
}

export interface WSGameEnd {
  sessionId: string
  totalScore: number
  rank: number
}

