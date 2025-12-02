import { useEffect, useCallback, useState } from 'react'
import { wsService } from '../services/websocket'
import type { Flight, WSMessage } from '../types'

export function useWebSocket(sessionId?: string) {
  const [isConnected, setIsConnected] = useState(false)
  const [flights, setFlights] = useState<Flight[]>([])

  useEffect(() => {
    wsService.connect(sessionId)
    setIsConnected(true)

    const unsubscribe = wsService.onFlightUpdate((updatedFlights) => {
      setFlights(updatedFlights)
    })

    return () => {
      unsubscribe()
    }
  }, [sessionId])

  const subscribe = useCallback((handler: (message: WSMessage) => void) => {
    return wsService.subscribe(handler)
  }, [])

  const send = useCallback((message: WSMessage) => {
    wsService.send(message)
  }, [])

  return {
    isConnected,
    flights,
    subscribe,
    send,
  }
}

