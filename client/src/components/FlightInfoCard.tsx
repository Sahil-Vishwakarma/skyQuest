import type { Flight } from '../types'

interface FlightInfoCardProps {
  flight: Flight | null
}

export function FlightInfoCard({ flight }: FlightInfoCardProps) {
  // Hint comes from server
  const hint = flight?.hint

  if (!flight) {
    return (
      <div className="card-glow p-4 animate-pulse">
        <div className="h-4 bg-midnight-700 rounded w-1/2 mb-3"></div>
        <div className="h-3 bg-midnight-700 rounded w-3/4 mb-2"></div>
        <div className="h-3 bg-midnight-700 rounded w-1/2"></div>
      </div>
    )
  }

  return (
    <div className="card-glow p-4">
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-sky-500/20 flex items-center justify-center">
            <svg className="w-4 h-4 text-sky-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/>
            </svg>
          </div>
          <div>
            {flight.flightNumber ? (
              <p className="font-mono font-bold text-white">{flight.flightNumber}</p>
            ) : (
              <p className="font-mono font-bold text-midnight-400">Flight ???</p>
            )}
          </div>
        </div>
      </div>

      {/* Airline */}
      {flight.airline.name && (
        <div className="mb-4">
          <p className="text-sm text-midnight-400 mb-1">Airline</p>
          <p className="text-white font-medium">{flight.airline.name}</p>
        </div>
      )}

      {/* Route */}
      <div className="flex items-center gap-3 mb-4">
        <div className="flex-1">
          <p className="text-xs text-midnight-400 mb-1">From</p>
          {flight.departure.iata !== '???' ? (
            <>
              <p className="font-mono font-bold text-white text-lg">{flight.departure.iata}</p>
              <p className="text-xs text-midnight-400">{flight.departure.city || flight.departure.name}</p>
            </>
          ) : (
            <>
              <p className="font-mono font-bold text-midnight-500 text-lg">???</p>
              <p className="text-xs text-midnight-500">Hidden</p>
            </>
          )}
        </div>
        
        <div className="flex-shrink-0">
          <svg className="w-6 h-6 text-midnight-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </div>
        
        <div className="flex-1 text-right">
          <p className="text-xs text-midnight-400 mb-1">To</p>
          <p className="font-mono font-bold text-amber-400 text-lg">???</p>
          <p className="text-xs text-amber-400/70">Your guess!</p>
        </div>
      </div>

      {/* Flight Stats */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-midnight-800/50 rounded-lg p-2 text-center">
          <p className="text-xs text-midnight-400 mb-1">Altitude</p>
          <p className="font-mono text-sm text-white">
            {Math.round(flight.altitude).toLocaleString()}
            <span className="text-midnight-400 text-xs"> ft</span>
          </p>
        </div>
        <div className="bg-midnight-800/50 rounded-lg p-2 text-center">
          <p className="text-xs text-midnight-400 mb-1">Speed</p>
          <p className="font-mono text-sm text-white">
            {Math.round(flight.speed)}
            <span className="text-midnight-400 text-xs"> kts</span>
          </p>
        </div>
        <div className="bg-midnight-800/50 rounded-lg p-2 text-center">
          <p className="text-xs text-midnight-400 mb-1">Heading</p>
          <p className="font-mono text-sm text-white">
            {Math.round(flight.direction)}°
          </p>
        </div>
      </div>

      {/* Aircraft */}
      {flight.aircraft.model && (
        <div className="mt-4 pt-4 border-t border-midnight-700/50">
          <p className="text-xs text-midnight-400 mb-1">Aircraft</p>
          <p className="text-sm text-white">{flight.aircraft.model}</p>
        </div>
      )}

      {/* Destination Hint */}
      {hint && (
        <div className="mt-4 pt-4 border-t border-midnight-700/50">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-medium text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-full border border-emerald-500/30">
              💡 Destination Hint
            </span>
          </div>
          <div className="bg-gradient-to-br from-emerald-900/20 to-sky-900/20 rounded-xl p-4 border border-emerald-500/20">
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-10 h-10 rounded-full bg-emerald-500/20 flex items-center justify-center">
                <span className="text-xl">🔍</span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm text-midnight-200 leading-relaxed">
                  <span className="text-amber-400">✨</span> {hint}
                </p>
                <p className="text-xs text-midnight-500 mt-2 italic">
                  Can you guess the city?
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
