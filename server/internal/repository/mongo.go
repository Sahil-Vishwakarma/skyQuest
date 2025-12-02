package repository

import (
	"context"
	"log"
	"time"

	"github.com/skyquest/server/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client   *mongo.Client
	db       *mongo.Database
	sessions *mongo.Collection
	users    *mongo.Collection
	scores   *mongo.Collection
}

func NewMongoRepository(uri, dbName string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	log.Println("MongoDB connected successfully")

	db := client.Database(dbName)

	repo := &MongoRepository{
		client:   client,
		db:       db,
		sessions: db.Collection("game_sessions"),
		users:    db.Collection("users"),
		scores:   db.Collection("leaderboard"),
	}

	// Create indexes
	if err := repo.createIndexes(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *MongoRepository) createIndexes(ctx context.Context) error {
	// Sessions collection indexes
	_, err := r.sessions.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "sessionId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "difficulty", Value: 1}}},
		{Keys: bson.D{{Key: "startedAt", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// Leaderboard collection indexes
	_, err = r.scores.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "difficulty", Value: 1}, {Key: "totalScore", Value: -1}}},
		{Keys: bson.D{{Key: "username", Value: 1}, {Key: "difficulty", Value: 1}}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

// Game Session methods

func (r *MongoRepository) CreateSession(ctx context.Context, session *models.GameSession) error {
	session.ID = primitive.NewObjectID()
	_, err := r.sessions.InsertOne(ctx, session)
	return err
}

func (r *MongoRepository) GetSession(ctx context.Context, sessionID string) (*models.GameSession, error) {
	var session models.GameSession
	err := r.sessions.FindOne(ctx, bson.M{"sessionId": sessionID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *MongoRepository) UpdateSession(ctx context.Context, session *models.GameSession) error {
	_, err := r.sessions.ReplaceOne(
		ctx,
		bson.M{"sessionId": session.SessionID},
		session,
	)
	return err
}

func (r *MongoRepository) GetUserSessions(ctx context.Context, userID string, limit int) ([]models.GameSession, error) {
	opts := options.Find().SetSort(bson.D{{Key: "startedAt", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.sessions.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.GameSession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// Leaderboard methods

func (r *MongoRepository) SaveScore(ctx context.Context, session *models.GameSession) error {
	filter := bson.M{
		"username":   session.Username,
		"difficulty": session.Difficulty,
	}

	// Check if user already has a score for this difficulty
	var existing models.LeaderboardEntry
	err := r.scores.FindOne(ctx, filter).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// Create new entry
		entry := models.LeaderboardEntry{
			ID:          primitive.NewObjectID(),
			Username:    session.Username,
			Difficulty:  session.Difficulty,
			TotalScore:  session.TotalScore,
			GamesPlayed: 1,
			UpdatedAt:   time.Now(),
		}
		_, err = r.scores.InsertOne(ctx, entry)
		return err
	}

	if err != nil {
		return err
	}

	// Update if new score is higher
	update := bson.M{
		"$inc": bson.M{"gamesPlayed": 1},
		"$set": bson.M{"updatedAt": time.Now()},
	}
	if session.TotalScore > existing.TotalScore {
		update["$set"].(bson.M)["totalScore"] = session.TotalScore
	}

	_, err = r.scores.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) GetLeaderboard(ctx context.Context, difficulty models.Difficulty, limit int) ([]models.LeaderboardEntry, error) {
	opts := options.Find().SetSort(bson.D{{Key: "totalScore", Value: -1}}).SetLimit(int64(limit))

	filter := bson.M{}
	if difficulty != "" {
		filter["difficulty"] = difficulty
	}

	cursor, err := r.scores.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []models.LeaderboardEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return entries, nil
}

func (r *MongoRepository) GetUserRank(ctx context.Context, username string, difficulty models.Difficulty) (int, error) {
	// Get user's score
	var userEntry models.LeaderboardEntry
	err := r.scores.FindOne(ctx, bson.M{
		"username":   username,
		"difficulty": difficulty,
	}).Decode(&userEntry)
	if err != nil {
		return 0, err
	}

	// Count users with higher scores
	count, err := r.scores.CountDocuments(ctx, bson.M{
		"difficulty": difficulty,
		"totalScore": bson.M{"$gt": userEntry.TotalScore},
	})
	if err != nil {
		return 0, err
	}

	return int(count) + 1, nil
}
