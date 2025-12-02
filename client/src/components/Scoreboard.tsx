import { formatScore } from '../utils/scoring'

interface ScoreboardProps {
  score: number
  round: number
  totalRounds: number
}

export function Scoreboard({ score, round, totalRounds }: ScoreboardProps) {
  const progress = (round / totalRounds) * 100

  return (
    <div className="flex items-center gap-4">
      {/* Progress */}
      <div className="hidden sm:flex items-center gap-2">
        <div className="w-24 h-2 bg-midnight-700 rounded-full overflow-hidden">
          <div 
            className="h-full bg-gradient-to-r from-sky-500 to-sky-400 transition-all duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
        <span className="text-xs text-midnight-400 font-mono">
          {round}/{totalRounds}
        </span>
      </div>

      {/* Score */}
      <div className="flex items-center gap-2 bg-midnight-800 px-4 py-2 rounded-xl border border-midnight-600">
        <svg className="w-4 h-4 text-amber-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
        <span className="font-mono font-bold text-white">
          {formatScore(score)}
        </span>
      </div>
    </div>
  )
}
