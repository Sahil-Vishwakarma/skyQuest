import type { Difficulty } from '../types'

export function getDifficultyLabel(difficulty: Difficulty): string {
  const labels: Record<Difficulty, string> = {
    easy: 'Easy',
    medium: 'Medium',
    hard: 'Hard',
  }
  return labels[difficulty]
}

export function getDifficultyColor(difficulty: Difficulty): string {
  const colors: Record<Difficulty, string> = {
    easy: 'text-emerald-400',
    medium: 'text-amber-400',
    hard: 'text-rose-400',
  }
  return colors[difficulty]
}

export function getDifficultyBgColor(difficulty: Difficulty): string {
  const colors: Record<Difficulty, string> = {
    easy: 'bg-emerald-500/20 border-emerald-500/30',
    medium: 'bg-amber-500/20 border-amber-500/30',
    hard: 'bg-rose-500/20 border-rose-500/30',
  }
  return colors[difficulty]
}

export function getMatchTypeLabel(matchType: string): string {
  const labels: Record<string, string> = {
    exact: 'Perfect Match!',
    family: 'Same City',
    country: 'Same Country',
    distance: 'Close Guess',
    wrong: 'Wrong',
  }
  return labels[matchType] || 'Unknown'
}

export function getMatchTypeColor(matchType: string): string {
  const colors: Record<string, string> = {
    exact: 'text-emerald-400',
    family: 'text-sky-400',
    country: 'text-amber-400',
    distance: 'text-orange-400',
    wrong: 'text-rose-400',
  }
  return colors[matchType] || 'text-gray-400'
}

export function formatScore(score: number): string {
  return score.toLocaleString()
}

export function formatTime(seconds: number): string {
  if (seconds < 60) {
    return `${seconds.toFixed(1)}s`
  }
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

export function formatDistance(km: number): string {
  if (km < 1) {
    return `${Math.round(km * 1000)}m`
  }
  if (km < 100) {
    return `${km.toFixed(1)}km`
  }
  return `${Math.round(km).toLocaleString()}km`
}

export function getSpeedBonusLabel(multiplier: number): string {
  if (multiplier >= 1.3) return 'Lightning Fast!'
  if (multiplier >= 1.1) return 'Quick!'
  return ''
}

export function calculateMaxScore(difficulty: Difficulty): number {
  const multipliers: Record<Difficulty, number> = {
    easy: 1.0,
    medium: 1.5,
    hard: 2.0,
  }
  // Max: 1000 base * difficulty * 1.3 speed bonus * 10 rounds
  return Math.round(1000 * multipliers[difficulty] * 1.3 * 10)
}

