import { useGameStore } from './store/gameStore'
import { Lobby } from './components/Lobby'
import { GameScreen } from './components/GameScreen'
import { GameOver } from './components/GameOver'

function App() {
  const status = useGameStore((state) => state.status)

  return (
    <div className="min-h-screen">
      {status === 'idle' && <Lobby />}
      {status === 'playing' && <GameScreen />}
      {status === 'finished' && <GameOver />}
    </div>
  )
}

export default App

