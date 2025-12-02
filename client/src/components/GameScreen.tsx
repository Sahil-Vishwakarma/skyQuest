import { useState } from 'react'
import { useGame } from '../hooks/useGame'
import { MapComponent } from './MapComponent'
import { FlightInfoCard } from './FlightInfoCard'
import { GuessInterface } from './GuessInterface'
import { ResultCard } from './ResultCard'
import { Scoreboard } from './Scoreboard'

export function GameScreen() {
  const { 
    currentRound, 
    totalRounds, 
    currentFlight, 
    totalScore, 
    showResult, 
    lastScore,
    nextRound,
    selectedAirport,
    selectAirport,
    submitGuess,
    isSubmitting,
    endGame,
    isEnding,
  } = useGame()

  const [showConfirmEnd, setShowConfirmEnd] = useState(false)

  const handleGuess = () => {
    if (selectedAirport) {
      submitGuess(selectedAirport)
    }
  }

  const handleEndGame = () => {
    setShowConfirmEnd(false)
    endGame() // Call the API endpoint to end game, then reset to lobby
  }

  return (
    <div className="min-h-screen flex flex-col">
      {/* Confirmation Modal */}
      {showConfirmEnd && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-midnight-800 border border-midnight-600 rounded-2xl p-6 max-w-md w-full animate-scale-in">
            <h3 className="text-xl font-bold text-white mb-2">End Game?</h3>
            <p className="text-midnight-300 mb-6">
              Are you sure you want to end this game? You'll return to the home screen and your progress will be lost.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowConfirmEnd(false)}
                className="flex-1 px-4 py-2 rounded-xl border border-midnight-500 text-midnight-300 hover:bg-midnight-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleEndGame}
                disabled={isEnding}
                className="flex-1 px-4 py-2 rounded-xl bg-rose-500 text-white hover:bg-rose-600 transition-colors disabled:opacity-50"
              >
                {isEnding ? 'Ending...' : 'End Game'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Header */}
      <header className="bg-midnight-900/80 backdrop-blur-xl border-b border-midnight-700/50 px-4 py-3">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-sky-500 to-indigo-600 flex items-center justify-center">
              <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
              </svg>
            </div>
            <div>
              <h1 className="font-display font-bold text-white">SkyQuest</h1>
              <p className="text-xs text-midnight-400">
                Round {currentRound} of {totalRounds}
              </p>
            </div>
          </div>
          
          <div className="flex items-center gap-3">
            <Scoreboard score={totalScore} round={currentRound} totalRounds={totalRounds} />
            
            {/* End Game Button */}
            <button
              onClick={() => setShowConfirmEnd(true)}
              className="px-3 py-1.5 rounded-lg bg-rose-500/20 border border-rose-500/30 text-rose-400 hover:bg-rose-500/30 transition-colors text-sm font-medium ml-4"
              title="End Game"
            >
              End Game
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex flex-col lg:flex-row">
        {/* Map Section */}
        <div className="flex-1 relative">
          <MapComponent 
            flight={currentFlight} 
            onAirportSelect={selectAirport}
            selectedAirport={selectedAirport}
          />
        </div>

        {/* Sidebar */}
        <div className="w-full lg:w-96 bg-midnight-900/95 backdrop-blur-xl border-l border-midnight-700/50 flex flex-col">
          {/* Flight Info */}
          <div className="p-4 border-b border-midnight-700/50">
            <FlightInfoCard flight={currentFlight} />
          </div>

          {/* Guess or Result */}
          <div className="flex-1 p-4 overflow-y-auto">
            {showResult && lastScore ? (
              <ResultCard 
                score={lastScore} 
                onNext={nextRound}
                roundNumber={currentRound}
                totalRounds={totalRounds}
              />
            ) : (
              <GuessInterface
                selectedAirport={selectedAirport}
                onSelect={selectAirport}
                onSubmit={handleGuess}
                isSubmitting={isSubmitting}
              />
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
