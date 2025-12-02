import { useState, useEffect, useRef } from 'react'
import { searchAirports, getAirportByIata } from '../utils/airports'
import type { Airport } from '../types'
import { cn } from '../utils/cn'

interface GuessInterfaceProps {
  selectedAirport: string | null
  onSelect: (iata: string | null) => void
  onSubmit: () => void
  isSubmitting: boolean
}

export function GuessInterface({ selectedAirport, onSelect, onSubmit, isSubmitting }: GuessInterfaceProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [suggestions, setSuggestions] = useState<Airport[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const suggestionsRef = useRef<HTMLDivElement>(null)

  const selectedAirportData = selectedAirport ? getAirportByIata(selectedAirport) : null

  useEffect(() => {
    if (searchQuery.length >= 2) {
      const results = searchAirports(searchQuery, 8)
      setSuggestions(results)
      setShowSuggestions(results.length > 0)
    } else {
      setSuggestions([])
      setShowSuggestions(false)
    }
  }, [searchQuery])

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(event.target as Node) &&
        !inputRef.current?.contains(event.target as Node)
      ) {
        setShowSuggestions(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleSelectAirport = (airport: Airport) => {
    onSelect(airport.iata)
    setSearchQuery('')
    setShowSuggestions(false)
  }

  const handleClear = () => {
    onSelect(null)
    setSearchQuery('')
  }

  return (
    <div className="space-y-4">
      <div>
        <h3 className="font-display font-semibold text-white mb-1">
          Where is this flight heading?
        </h3>
        <p className="text-sm text-midnight-400">
          Search for an airport or click on the map
        </p>
      </div>

      {/* Search Input */}
      <div className="relative">
        <div className="relative">
          <input
            ref={inputRef}
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onFocus={() => searchQuery.length >= 2 && setShowSuggestions(true)}
            placeholder="Search airport (e.g., LAX, London, Tokyo)..."
            className="input-field pr-10"
          />
          <svg
            className="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-midnight-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>

        {/* Suggestions Dropdown */}
        {showSuggestions && suggestions.length > 0 && (
          <div
            ref={suggestionsRef}
            className="absolute z-50 w-full mt-2 bg-midnight-800 border border-midnight-600 rounded-xl shadow-xl overflow-hidden"
          >
            {suggestions.map((airport) => (
              <button
                key={airport.iata}
                onClick={() => handleSelectAirport(airport)}
                className={cn(
                  'w-full px-4 py-3 text-left hover:bg-midnight-700 transition-colors flex items-center gap-3',
                  selectedAirport === airport.iata && 'bg-sky-500/10'
                )}
              >
                <span className="font-mono font-bold text-sky-400 w-12">
                  {airport.iata}
                </span>
                <div className="flex-1 min-w-0">
                  <p className="text-white truncate">{airport.name}</p>
                  <p className="text-xs text-midnight-400 truncate">
                    {airport.city}, {airport.country}
                  </p>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Selected Airport Display */}
      {selectedAirportData && (
        <div className="bg-sky-500/10 border border-sky-500/30 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-sky-500/20 flex items-center justify-center">
                <span className="font-mono font-bold text-sky-400">
                  {selectedAirportData.iata}
                </span>
              </div>
              <div>
                <p className="text-white font-medium">{selectedAirportData.name}</p>
                <p className="text-sm text-midnight-400">
                  {selectedAirportData.city}, {selectedAirportData.country}
                </p>
              </div>
            </div>
            <button
              onClick={handleClear}
              className="p-2 hover:bg-midnight-700 rounded-lg transition-colors"
            >
              <svg className="w-5 h-5 text-midnight-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* Submit Button */}
      <button
        onClick={onSubmit}
        disabled={!selectedAirport || isSubmitting}
        className="btn-primary w-full"
      >
        {isSubmitting ? (
          <span className="flex items-center justify-center gap-2">
            <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Checking...
          </span>
        ) : selectedAirport ? (
          `Guess: ${selectedAirport}`
        ) : (
          'Select a destination'
        )}
      </button>

      {/* Hint */}
      <p className="text-xs text-midnight-500 text-center">
        Tip: Click on any airport marker on the map to select it
      </p>
    </div>
  )
}

