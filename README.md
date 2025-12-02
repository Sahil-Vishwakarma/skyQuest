# SkyQuest - Flight Guessing Game

An interactive web-based game that leverages real-time flight tracking data to create an engaging experience. Players observe live aircraft on a map, make educated guesses about destination airports, and earn points based on accuracy.

## Tech Stack

### Frontend
- React 18 + TypeScript
- Vite (build tool)
- Leaflet.js (interactive maps)
- TailwindCSS (styling)

### Backend
- Go 1.21+
- Gin (HTTP framework)
- gorilla/websocket (real-time updates)
- MongoDB (database)
- Redis (caching)

### External APIs
- Aviation Edge ADSB API (flight tracking)

## Project Structure

```
SkyQuest/
├── client/                 # React frontend
│   ├── src/
│   │   ├── components/     # UI components
│   │   ├── hooks/          # Custom React hooks
│   │   ├── services/       # API clients
│   │   ├── types/          # TypeScript types
│   │   └── utils/          # Utility functions
│   └── package.json
├── server/                 # Go backend
│   ├── cmd/api/            # Entry point
│   ├── internal/           # Private packages
│   │   ├── config/         # Configuration
│   │   ├── handlers/       # HTTP handlers
│   │   ├── models/         # Data models
│   │   ├── repository/     # Data access
│   │   ├── services/       # Business logic
│   │   └── websocket/      # WebSocket hub
│   └── pkg/aviation/       # Aviation API client
├── docker-compose.yml      # MongoDB + Redis
└── README.md
```

## Getting Started

### Prerequisites
- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- Aviation Edge API key (free tier available)

### 1. Start Infrastructure

```bash
docker-compose up -d
```

This starts MongoDB and Redis.

### 2. Configure Environment

```bash
# Server
cd server
cp .env.example .env
# Edit .env with your Aviation Edge API key
```

### 3. Run Backend

```bash
cd server
go mod download
go run cmd/api/main.go
```

Server runs on http://localhost:8080

### 4. Run Frontend

```bash
cd client
npm install
npm run dev
```

Frontend runs on http://localhost:5173

## Game Rules

### Scoring
| Result | Points |
|--------|--------|
| Exact Match | 1000 |
| Airport Family (same city) | 750 |
| Correct Region | 500 |
| Within 500km | 250 |
| Wrong | 0 |

### Difficulty Multipliers
- Easy: 1.0x (domestic flights)
- Medium: 1.5x (international, same region)
- Hard: 2.0x (global flights)

### Speed Bonuses
- Under 10 seconds: 1.3x
- Under 30 seconds: 1.1x
- Over 30 seconds: 1.0x

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/flights` | Get available flights |
| POST | `/api/game/start` | Start new game |
| POST | `/api/game/guess` | Submit guess |
| POST | `/api/game/end` | End game |
| GET | `/api/leaderboard` | Get leaderboard |
| WS | `/ws` | WebSocket connection |

## Development

### Hot Reload (Go)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd server
air
```

### Linting
```bash
# Go
cd server
go vet ./...

# Frontend
cd client
npm run lint
```

## License

MIT

# skyQuest
