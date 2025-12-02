package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/skyquest/server/internal/config"
	"github.com/skyquest/server/internal/handlers"
	"github.com/skyquest/server/internal/repository"
	"github.com/skyquest/server/internal/services"
	"github.com/skyquest/server/internal/websocket"
	"github.com/skyquest/server/pkg/aviation"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	mongoRepo, err := repository.NewMongoRepository(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to MongoDB: %v", err)
		log.Println("Game sessions and leaderboard will not be persisted.")
		log.Println("Start MongoDB with: docker-compose up -d")
	}
	if mongoRepo != nil {
		defer mongoRepo.Close()
	}

	// Initialize Redis
	redisClient := repository.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if redisClient != nil {
		defer redisClient.Close()
	}

	// Initialize AviationStack API client
	aviationClient := aviation.NewClient(cfg.AviationStackAPIKey)

	// Initialize services
	flightService := services.NewFlightService(aviationClient, redisClient)
	var gameService *services.GameService
	var scoreService *services.ScoreService
	if mongoRepo != nil {
		gameService = services.NewGameService(mongoRepo, flightService)
		scoreService = services.NewScoreService(mongoRepo)
	} else {
		// Create services with nil repo (limited functionality)
		gameService = services.NewGameService(nil, flightService)
		scoreService = services.NewScoreService(nil)
	}

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Start flight data polling in background
	go flightService.StartPolling(wsHub, 5*time.Minute)

	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(gameService, scoreService)
	flightHandler := handlers.NewFlightHandler(flightService)
	leaderboardHandler := handlers.NewLeaderboardHandler(scoreService)
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// API routes
	api := router.Group("/api")
	{
		// Flight endpoints
		api.GET("/flights", flightHandler.GetFlights)

		// Game endpoints
		api.POST("/game/start", gameHandler.StartGame)
		api.POST("/game/guess", gameHandler.SubmitGuess)
		api.POST("/game/end", gameHandler.EndGame)

		// Leaderboard endpoints
		api.GET("/leaderboard", leaderboardHandler.GetLeaderboard)
	}

	// WebSocket endpoint
	router.GET("/ws", wsHandler.HandleWebSocket)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
