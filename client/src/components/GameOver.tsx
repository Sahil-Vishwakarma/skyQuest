import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useGameStore } from '../store/gameStore'
import { api } from '../services/api'
import { cn } from '../utils/cn'
import { formatScore } from '../utils/scoring'
import type { Round } from '../types'

export function GameOver() {
  const { 
    totalScore, 
    username, 
    resetGame,
    sessionId,
  } = useGameStore()
  
  const [showConfetti, setShowConfetti] = useState(true)
  const [scoreSaved, setScoreSaved] = useState(false)
  const [rounds, setRounds] = useState<Round[]>([])

  // End game to save score first and get rounds data
  useEffect(() => {
    if (sessionId && !scoreSaved) {
      api.endGame(sessionId)
        .then((response) => {
          setRounds(response.rounds || [])
          setScoreSaved(true)
        })
        .catch(console.error)
    }
  }, [sessionId, scoreSaved])

  // Fetch leaderboard after score is saved
  const { data: leaderboardData, isLoading: leaderboardLoading } = useQuery({
    queryKey: ['leaderboard', scoreSaved],
    queryFn: () => api.getLeaderboard('easy', 10),
    enabled: scoreSaved,
  })

  useEffect(() => {
    const timer = setTimeout(() => setShowConfetti(false), 3000)
    return () => clearTimeout(timer)
  }, [])

  const maxScore = 10000 // 10 rounds * 1000 max per round
  const scorePercentage = Math.round((totalScore / maxScore) * 100)
  const totalRounds = rounds.length || 10
  const correctGuesses = rounds.filter(r => r.pointsEarned >= 1000).length
  const partialGuesses = rounds.filter(r => r.pointsEarned > 0 && r.pointsEarned < 1000).length
  const missedGuesses = totalRounds - correctGuesses - partialGuesses

  const getRank = () => {
    if (scorePercentage >= 90) return { title: 'Aviation Expert', emoji: '🏆' }
    if (scorePercentage >= 70) return { title: 'Skilled Navigator', emoji: '🥇' }
    if (scorePercentage >= 50) return { title: 'Frequent Flyer', emoji: '🥈' }
    if (scorePercentage >= 30) return { title: 'Amateur Spotter', emoji: '🥉' }
    return { title: 'Ground Crew', emoji: '✈️' }
  }

  const rank = getRank()

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        {/* Confetti Effect */}
        {showConfetti && scorePercentage >= 50 && (
          <div className="fixed inset-0 pointer-events-none overflow-hidden">
            {[...Array(50)].map((_, i) => (
              <div
                key={i}
                className="absolute animate-float"
                style={{
                  left: `${Math.random() * 100}%`,
                  top: `-${Math.random() * 20}%`,
                  animationDelay: `${Math.random() * 2}s`,
                  animationDuration: `${3 + Math.random() * 2}s`,
                }}
              >
                <span className="text-2xl">
                  {['🎉', '✨', '⭐', '🎊', '🏆'][Math.floor(Math.random() * 5)]}
                </span>
              </div>
            ))}
          </div>
        )}

        {/* Main Card */}
        <div className="card p-8 animate-scale-in">
          {/* Header */}
          <div className="text-center mb-8">
            <div className="text-6xl mb-4">{rank.emoji}</div>
            <h1 className="text-3xl font-display font-bold text-white mb-2">
              Game Over!
            </h1>
            <p className="text-xl text-sky-400">{rank.title}</p>
            <p className="text-midnight-400 mt-1">{username}</p>
          </div>

          {/* Score Display */}
          <div className="bg-gradient-to-br from-sky-500/20 to-indigo-500/20 rounded-2xl p-6 mb-6 text-center border border-sky-500/30">
            <p className="text-sm text-midnight-300 mb-1">Final Score</p>
            <p className="text-5xl font-display font-bold text-white mb-2">
              {formatScore(totalScore)}
            </p>
            <p className="text-midnight-400 text-sm">
              {scorePercentage}% of max ({formatScore(maxScore)})
            </p>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-3 gap-4 mb-6">
            <div className="bg-midnight-800/50 rounded-xl p-4 text-center">
              <p className="text-2xl font-bold text-emerald-400">
                {rounds.length > 0 ? correctGuesses : '—'}
              </p>
              <p className="text-xs text-midnight-400">Perfect</p>
            </div>
            <div className="bg-midnight-800/50 rounded-xl p-4 text-center">
              <p className="text-2xl font-bold text-amber-400">
                {rounds.length > 0 ? partialGuesses : '—'}
              </p>
              <p className="text-xs text-midnight-400">Partial</p>
            </div>
            <div className="bg-midnight-800/50 rounded-xl p-4 text-center">
              <p className="text-2xl font-bold text-rose-400">
                {rounds.length > 0 ? missedGuesses : '—'}
              </p>
              <p className="text-xs text-midnight-400">Missed</p>
            </div>
          </div>

          {/* Round Breakdown */}
          <div className="mb-6">
            <h3 className="text-sm font-medium text-midnight-300 mb-3">Round Breakdown</h3>
            {rounds.length === 0 ? (
              <div className="text-center py-4">
                <div className="animate-spin h-5 w-5 border-2 border-sky-500 border-t-transparent rounded-full mx-auto mb-2"></div>
                <p className="text-midnight-500 text-sm">Loading rounds...</p>
              </div>
            ) : (
              <div className="space-y-2 max-h-48 overflow-y-auto">
                {rounds.map((round, index) => (
                  <div 
                    key={index}
                    className="flex items-center justify-between bg-midnight-800/50 rounded-lg px-3 py-2"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-midnight-500 text-sm w-6">#{round.roundNumber}</span>
                      <span className="font-mono text-white">{round.playerGuess || '—'}</span>
                      <span className="text-midnight-500">→</span>
                      <span className="font-mono text-emerald-400">{round.actualArrival}</span>
                    </div>
                    <span className={cn(
                      'font-mono text-sm',
                      round.pointsEarned >= 1000 ? 'text-emerald-400' :
                      round.pointsEarned > 0 ? 'text-amber-400' : 'text-midnight-500'
                    )}>
                      +{formatScore(round.pointsEarned)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Leaderboard Preview */}
          <div className="mb-6">
            <h3 className="text-sm font-medium text-midnight-300 mb-3">Leaderboard</h3>
            {leaderboardLoading ? (
              <div className="text-center py-4">
                <div className="animate-spin h-6 w-6 border-2 border-sky-500 border-t-transparent rounded-full mx-auto mb-2"></div>
                <p className="text-midnight-500 text-sm">Loading leaderboard...</p>
              </div>
            ) : leaderboardData && leaderboardData.leaderboard && leaderboardData.leaderboard.length > 0 ? (
              <div className="space-y-2">
                {leaderboardData.leaderboard.slice(0, 5).map((entry, index) => (
                  <div 
                    key={entry.id || index}
                    className={cn(
                      'flex items-center justify-between rounded-lg px-3 py-2',
                      entry.username === username 
                        ? 'bg-sky-500/10 border border-sky-500/30' 
                        : 'bg-midnight-800/50'
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <span className={cn(
                        'w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold',
                        index === 0 && 'bg-amber-500 text-midnight-900',
                        index === 1 && 'bg-midnight-400 text-midnight-900',
                        index === 2 && 'bg-amber-700 text-white',
                        index > 2 && 'bg-midnight-700 text-midnight-300'
                      )}>
                        {index + 1}
                      </span>
                      <span className={cn(
                        'font-medium',
                        entry.username === username ? 'text-sky-400' : 'text-white'
                      )}>
                        {entry.username}
                      </span>
                    </div>
                    <span className="font-mono text-midnight-300">
                      {formatScore(entry.totalScore)}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="bg-midnight-800/50 rounded-lg px-4 py-6 text-center">
                <p className="text-midnight-400 text-sm mb-1">No scores yet</p>
                <p className="text-midnight-500 text-xs">
                  {!scoreSaved ? 'Saving your score...' : 'Be the first on the leaderboard!'}
                </p>
              </div>
            )}
          </div>

          {/* Actions */}
          <div className="flex gap-3">
            <button
              onClick={resetGame}
              className="btn-primary flex-1"
            >
              Play Again
            </button>
            <button
              onClick={() => {
                const text = `🛫 I scored ${formatScore(totalScore)} points in SkyQuest!\n` +
                  `Rank: ${rank.title} ${rank.emoji}\n` +
                  `Accuracy: ${correctGuesses}/10 perfect guesses`
                navigator.clipboard.writeText(text)
              }}
              className="btn-secondary px-4"
              title="Copy results"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
              </svg>
            </button>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center mt-6 text-midnight-500 text-sm">
          Thanks for playing SkyQuest!
        </div>
      </div>
    </div>
  )
}
