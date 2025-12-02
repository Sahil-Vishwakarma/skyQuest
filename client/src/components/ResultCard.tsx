import type { ScoreResult } from '../types'
import { cn } from '../utils/cn'
import { 
  getMatchTypeLabel, 
  getMatchTypeColor, 
  formatScore, 
  formatDistance,
  getSpeedBonusLabel,
} from '../utils/scoring'

interface ResultCardProps {
  score: ScoreResult
  onNext: () => void
  roundNumber: number
  totalRounds: number
}

export function ResultCard({ score, onNext, roundNumber, totalRounds }: ResultCardProps) {
  const isLastRound = roundNumber === totalRounds
  const isCorrect = score.matchType === 'exact'
  const isPartialCorrect = ['family', 'country', 'distance'].includes(score.matchType)
  const speedBonus = getSpeedBonusLabel(score.speedMultiplier)

  return (
    <div className="space-y-4 animate-scale-in">
      {/* Result Header */}
      <div className={cn(
        'text-center p-6 rounded-xl border',
        isCorrect && 'bg-emerald-500/10 border-emerald-500/30',
        isPartialCorrect && 'bg-amber-500/10 border-amber-500/30',
        !isCorrect && !isPartialCorrect && 'bg-rose-500/10 border-rose-500/30',
      )}>
        <div className={cn(
          'text-4xl mb-2',
          isCorrect && 'text-emerald-400',
          isPartialCorrect && 'text-amber-400',
          !isCorrect && !isPartialCorrect && 'text-rose-400',
        )}>
          {isCorrect ? '🎯' : isPartialCorrect ? '📍' : '❌'}
        </div>
        <h3 className={cn('text-xl font-display font-bold', getMatchTypeColor(score.matchType))}>
          {getMatchTypeLabel(score.matchType)}
        </h3>
        {speedBonus && (
          <p className="text-sky-400 text-sm mt-1">{speedBonus}</p>
        )}
      </div>

      {/* Airport Comparison */}
      <div className="bg-midnight-800/50 rounded-xl p-4">
        <div className="flex items-center justify-between">
          <div className="text-center flex-1">
            <p className="text-xs text-midnight-400 mb-1">Your Guess</p>
            <p className="font-mono font-bold text-lg text-white">
              {score.guessedAirport?.iata || '???'}
            </p>
            <p className="text-xs text-midnight-400">
              {score.guessedAirport?.city || 'Unknown'}
            </p>
          </div>
          
          <div className="px-4">
            <svg className="w-6 h-6 text-midnight-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
          </div>
          
          <div className="text-center flex-1">
            <p className="text-xs text-midnight-400 mb-1">Correct Answer</p>
            <p className="font-mono font-bold text-lg text-emerald-400">
              {score.correctAirport?.iata || '???'}
            </p>
            <p className="text-xs text-midnight-400">
              {score.correctAirport?.city || 'Unknown'}
            </p>
          </div>
        </div>

        {score.distanceKm > 0 && (
          <div className="mt-3 pt-3 border-t border-midnight-700 text-center">
            <p className="text-sm text-midnight-400">
              Distance: <span className="text-white font-medium">{formatDistance(score.distanceKm)}</span>
            </p>
          </div>
        )}
      </div>

      {/* Score Breakdown */}
      <div className="bg-midnight-800/50 rounded-xl p-4 space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-midnight-400">Base Points</span>
          <span className="text-white font-mono">{formatScore(score.basePoints)}</span>
        </div>
        {score.speedMultiplier > 1 && (
          <div className="flex justify-between text-sm">
            <span className="text-midnight-400">Speed Bonus</span>
            <span className="text-sky-400 font-mono">×{score.speedMultiplier}</span>
          </div>
        )}
        <div className="flex justify-between text-lg font-bold pt-2 border-t border-midnight-700">
          <span className="text-white">Total</span>
          <span className={cn(
            'font-mono',
            score.totalPoints > 0 ? 'text-emerald-400' : 'text-midnight-400'
          )}>
            +{formatScore(score.totalPoints)}
          </span>
        </div>
      </div>

      {/* Next Button */}
      <button onClick={onNext} className="btn-primary w-full">
        {isLastRound ? 'See Final Results' : `Next Round (${roundNumber + 1}/${totalRounds})`}
      </button>
    </div>
  )
}
