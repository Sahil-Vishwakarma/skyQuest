import { create } from 'zustand'
import type { 
  GameStatus, 
  Flight, 
  ScoreResult,
  Round,
} from '../types'

interface GameState {
  // Game status
  status: GameStatus
  sessionId: string | null
  username: string
  
  // Current round
  currentRound: number
  totalRounds: number
  currentFlight: Flight | null
  roundStartTime: number | null
  
  // Scoring
  totalScore: number
  rounds: Round[]
  lastScore: ScoreResult | null
  
  // UI state
  selectedAirport: string | null
  showResult: boolean
  
  // Actions
  setUsername: (username: string) => void
  startGame: (sessionId: string, flight: Flight, totalRounds: number) => void
  setCurrentFlight: (flight: Flight) => void
  selectAirport: (iata: string | null) => void
  submitGuess: (score: ScoreResult, nextFlight: Flight | null, isGameOver: boolean) => void
  nextRound: () => void
  endGame: (rounds: Round[], rank: number) => void
  resetGame: () => void
}

const initialState = {
  status: 'idle' as GameStatus,
  sessionId: null,
  username: '',
  currentRound: 0,
  totalRounds: 10,
  currentFlight: null,
  roundStartTime: null,
  totalScore: 0,
  rounds: [],
  lastScore: null,
  selectedAirport: null,
  showResult: false,
}

export const useGameStore = create<GameState>((set) => ({
  ...initialState,

  setUsername: (username) => set({ username }),
  
  startGame: (sessionId, flight, totalRounds) => set({
    status: 'playing',
    sessionId,
    currentRound: 1,
    totalRounds,
    currentFlight: flight,
    roundStartTime: Date.now(),
    totalScore: 0,
    rounds: [],
    lastScore: null,
    selectedAirport: null,
    showResult: false,
  }),
  
  setCurrentFlight: (flight) => set({
    currentFlight: flight,
    roundStartTime: Date.now(),
  }),
  
  selectAirport: (iata) => set({ selectedAirport: iata }),
  
  submitGuess: (score, nextFlight, isGameOver) => set((state) => ({
    lastScore: score,
    totalScore: state.totalScore + score.totalPoints,
    showResult: true,
    currentFlight: isGameOver ? state.currentFlight : nextFlight,
    status: isGameOver ? 'finished' : state.status,
  })),
  
  nextRound: () => set((state) => ({
    currentRound: state.currentRound + 1,
    roundStartTime: Date.now(),
    selectedAirport: null,
    showResult: false,
    lastScore: null,
  })),
  
  endGame: (rounds, _rank) => set({
    status: 'finished',
    rounds,
    showResult: false,
  }),
  
  resetGame: () => set(() => ({ ...initialState }), true),
}))
