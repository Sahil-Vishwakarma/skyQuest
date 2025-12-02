import { useMutation } from '@tanstack/react-query'
import { api } from '../services/api'
import { useGameStore } from '../store/gameStore'
import type { StartGameRequest, GuessRequest } from '../types'

export function useGame() {
  const store = useGameStore()

  const startGameMutation = useMutation({
    mutationFn: (data: StartGameRequest) => api.startGame(data),
    onSuccess: (response) => {
      store.startGame(response.sessionId, response.flight, response.totalRounds)
    },
  })

  const submitGuessMutation = useMutation({
    mutationFn: (data: GuessRequest) => api.submitGuess(data),
    onSuccess: (response) => {
      store.submitGuess(response.score, response.nextFlight || null, response.isGameOver)
    },
  })

  const endGameMutation = useMutation({
    mutationFn: (sessionId: string) => api.endGame(sessionId),
    onSuccess: () => {
      // Reset game to go back to lobby after ending
      store.resetGame()
    },
    onError: () => {
      // Even if API fails, reset to lobby
      store.resetGame()
    },
  })

  const startGame = () => {
    if (!store.username.trim()) return
    
    startGameMutation.mutate({
      username: store.username,
      difficulty: 'easy', // Always use easy mode
    })
  }

  const submitGuess = (airportIata: string, confidence?: number) => {
    if (!store.sessionId) return
    
    submitGuessMutation.mutate({
      sessionId: store.sessionId,
      airportIata,
      confidence,
    })
  }

  const endGame = () => {
    if (!store.sessionId) return
    endGameMutation.mutate(store.sessionId)
  }

  return {
    ...store,
    startGame,
    submitGuess,
    endGame,
    isStarting: startGameMutation.isPending,
    isSubmitting: submitGuessMutation.isPending,
    isEnding: endGameMutation.isPending,
    startError: startGameMutation.error,
    submitError: submitGuessMutation.error,
  }
}

