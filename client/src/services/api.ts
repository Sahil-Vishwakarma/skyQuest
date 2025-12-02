import type {
  StartGameRequest,
  StartGameResponse,
  GuessRequest,
  GuessResponse,
  EndGameResponse,
  FlightsResponse,
  LeaderboardResponse,
  Difficulty,
} from '../types'

const API_BASE = '/api'

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function fetchJson<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new ApiError(response.status, error.error || 'Request failed')
  }

  return response.json()
}

export const api = {
  // Flight endpoints
  async getFlights(difficulty?: Difficulty, limit?: number): Promise<FlightsResponse> {
    const params = new URLSearchParams()
    if (difficulty) params.set('difficulty', difficulty)
    if (limit) params.set('limit', limit.toString())
    
    return fetchJson<FlightsResponse>(`${API_BASE}/flights?${params}`)
  },

  // Game endpoints
  async startGame(data: StartGameRequest): Promise<StartGameResponse> {
    return fetchJson<StartGameResponse>(`${API_BASE}/game/start`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async submitGuess(data: GuessRequest): Promise<GuessResponse> {
    return fetchJson<GuessResponse>(`${API_BASE}/game/guess`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async endGame(sessionId: string): Promise<EndGameResponse> {
    return fetchJson<EndGameResponse>(`${API_BASE}/game/end`, {
      method: 'POST',
      body: JSON.stringify({ sessionId }),
    })
  },

  // Leaderboard endpoints
  async getLeaderboard(difficulty?: Difficulty, limit?: number): Promise<LeaderboardResponse> {
    const params = new URLSearchParams()
    if (difficulty) params.set('difficulty', difficulty)
    if (limit) params.set('limit', limit.toString())
    
    return fetchJson<LeaderboardResponse>(`${API_BASE}/leaderboard?${params}`)
  },
}

export { ApiError }

