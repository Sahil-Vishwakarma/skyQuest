import type { WSMessage, Flight } from '../types'

type MessageHandler = (message: WSMessage) => void

class WebSocketService {
  private ws: WebSocket | null = null
  private handlers: Set<MessageHandler> = new Set()
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private sessionId: string | null = null

  connect(sessionId?: string): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      if (sessionId) {
        this.sessionId = sessionId
        this.send({ type: 'register', payload: { sessionId } })
      }
      return
    }

    const wsUrl = `ws://${window.location.host}/ws${sessionId ? `?sessionId=${sessionId}` : ''}`
    this.sessionId = sessionId || null

    try {
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
      }

      this.ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data)
          this.handlers.forEach((handler) => handler(message))
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.attemptReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
      this.attemptReconnect()
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached')
      return
    }

    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)

    console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`)

    setTimeout(() => {
      this.connect(this.sessionId || undefined)
    }, delay)
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.sessionId = null
  }

  send(message: WSMessage): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    }
  }

  subscribe(handler: MessageHandler): () => void {
    this.handlers.add(handler)
    return () => this.handlers.delete(handler)
  }

  onFlightUpdate(callback: (flights: Flight[]) => void): () => void {
    return this.subscribe((message) => {
      if (message.type === 'flight:update') {
        const payload = message.payload as { flights: Flight[] }
        callback(payload.flights)
      }
    })
  }
}

export const wsService = new WebSocketService()

