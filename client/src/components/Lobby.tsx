import { useState } from 'react'
import { useGame } from '../hooks/useGame'
import { cn } from '../utils/cn'

export function Lobby() {
  const { username, setUsername, startGame, isStarting, startError } = useGame()
  const [usernameError, setUsernameError] = useState('')

  const handleStart = () => {
    if (!username.trim()) {
      setUsernameError('Please enter a username')
      return
    }
    if (username.length < 2) {
      setUsernameError('Username must be at least 2 characters')
      return
    }
    setUsernameError('')
    startGame()
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-lg">
        {/* Logo & Title */}
        <div className="text-center mb-12 animate-slide-down">
          <div className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-sky-500 to-indigo-600 shadow-lg shadow-sky-500/30 mb-6">
            <svg className="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
            </svg>
          </div>
          <h1 className="text-5xl font-display font-bold bg-gradient-to-r from-white via-sky-200 to-sky-400 bg-clip-text text-transparent mb-3">
            SkyQuest
          </h1>
          <p className="text-midnight-300 text-lg">
            Test your aviation knowledge by guessing flight destinations
          </p>
        </div>

        {/* Main Card */}
        <div className="card p-8 animate-scale-in" style={{ animationDelay: '0.1s' }}>
          {/* Username */}
          <div className="mb-8">
            <label className="block text-sm font-medium text-midnight-300 mb-2">
              Your Name
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value)
                setUsernameError('')
              }}
              placeholder="Enter your pilot name..."
              className={cn('input-field', usernameError && 'border-rose-500 focus:border-rose-500')}
              maxLength={20}
              onKeyDown={(e) => e.key === 'Enter' && handleStart()}
            />
            {usernameError && (
              <p className="text-rose-400 text-sm mt-2">{usernameError}</p>
            )}
          </div>

          {/* Game Info */}
          <div className="mb-8 p-4 bg-emerald-500/10 border border-emerald-500/30 rounded-xl">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 rounded-full bg-emerald-500/20 flex items-center justify-center">
                <span className="text-xl">✈️</span>
              </div>
              <div>
                <h3 className="font-semibold text-white">How to Play</h3>
                <p className="text-sm text-emerald-400">10 rounds • Guess the destination</p>
              </div>
            </div>
            <ul className="text-sm text-midnight-300 space-y-1">
              <li className="flex items-center gap-2">
                <span className="text-emerald-400">✓</span> See the flight info and departure airport
              </li>
              <li className="flex items-center gap-2">
                <span className="text-emerald-400">✓</span> Get hints about the destination city
              </li>
              <li className="flex items-center gap-2">
                <span className="text-emerald-400">✓</span> Guess the arrival airport to score points
              </li>
            </ul>
          </div>

          {/* Error Message */}
          {startError && (
            <div className="mb-6 p-4 bg-rose-500/10 border border-rose-500/30 rounded-xl text-rose-400 text-sm">
              {startError instanceof Error ? startError.message : 'Failed to start game'}
            </div>
          )}

          {/* Start Button */}
          <button
            onClick={handleStart}
            disabled={isStarting}
            className="btn-primary w-full text-lg py-4"
          >
            {isStarting ? (
              <span className="flex items-center justify-center gap-2">
                <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Starting...
              </span>
            ) : (
              'Start Game'
            )}
          </button>
        </div>

        {/* Footer */}
        <div className="mt-8 text-center animate-fade-in" style={{ animationDelay: '0.3s' }}>
          <p className="text-midnight-500 text-sm">
            Track real flights and learn about airports around the world!
          </p>
        </div>
      </div>
    </div>
  )
}
