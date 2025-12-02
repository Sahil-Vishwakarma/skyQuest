// Frontend configuration
// Set VITE_API_URL in .env or environment variables to override

const config = {
  // API base URL - defaults to relative path for same-origin deployment
  apiUrl: import.meta.env.VITE_API_URL || '',
  
  // WebSocket URL - defaults to current host
  wsUrl: import.meta.env.VITE_WS_URL || '',
}

// Helper to get the full API URL
export function getApiUrl(path: string): string {
  return `${config.apiUrl}${path}`
}

// Helper to get the WebSocket URL
export function getWsUrl(sessionId?: string): string {
  if (config.wsUrl) {
    // Use configured WebSocket URL
    return `${config.wsUrl}${sessionId ? `?sessionId=${sessionId}` : ''}`
  }
  
  // Default: derive from current location
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/ws${sessionId ? `?sessionId=${sessionId}` : ''}`
}

export default config

