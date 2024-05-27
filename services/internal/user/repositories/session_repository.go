package repositories

import (
	"context"
	"errors"
	"fmt"
	"services/internal/user/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionRepository struct {
	collection mongo.Collection
}

func NewSessionRepository(db *mongo.Database, collectionName string) (*SessionRepository, error) {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		collection: *db.Collection(collectionName),
	}, nil
}

func (r *SessionRepository) Create(Session *models.Session) (*models.Session, error) {
	Session.CreatedAt = time.Now()
	Session.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(context.TODO(), Session)
	if err != nil {
		return nil, err
	}
	Session.ID = res.InsertedID.(primitive.ObjectID)

	return Session, nil
}

// GetAll tüm kitapları getirir
func (r *SessionRepository) GetAll() ([]*models.Session, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var books []*models.Session
	for cur.Next(context.Background()) {
		var book models.Session
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (r *SessionRepository) Delete(SessionID primitive.ObjectID) error {
	filter := bson.M{"_id": SessionID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Session Not Found.")
	}
	return nil
}

// FindOne belirli bir hesabı ID ile bulur
func (r *SessionRepository) FindOne(SessionID primitive.ObjectID) (*models.Session, error) {
	filter := bson.M{"_id": SessionID}
	var Session models.Session
	err := r.collection.FindOne(context.TODO(), filter).Decode(&Session)
	if err != nil {
		return nil, err
	}
	return &Session, nil
}
func (r *SessionRepository) FindOneWithParameter(key string, data any) (*models.Session, error) {
	filter := bson.M{key: data}
	var Session models.Session
	err := r.collection.FindOne(context.TODO(), filter).Decode(&Session)
	if err != nil {
		return nil, err
	}
	return &Session, nil
}

// FindAll tüm hesapları getirir
func (r *SessionRepository) FindAll() ([]*models.Session, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var Sessions []*models.Session
	for cur.Next(context.Background()) {
		var Session models.Session
		err := cur.Decode(&Session)
		if err != nil {
			return nil, err
		}
		Sessions = append(Sessions, &Session)
	}

	return Sessions, nil
}

// Update bir hesabın bilgilerini günceller
func (r *SessionRepository) Update(SessionID primitive.ObjectID, updatedSession *models.Session) error {
	filter := bson.M{"_id": SessionID}
	updatedSession.UpdatedAt = time.Now()
	update := bson.M{
		"$set": updatedSession,
	}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Session Not Found.")
	}
	return nil
}

func (r *SessionRepository) UpdateOrCreateSession(session *models.Session, keyType string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M
	switch keyType {
	case "access_token":
		filter = bson.M{"access_token": session.AccessToken}
	case "refresh_token":
		filter = bson.M{"refresh_token": session.RefreshToken}
	case "account":
		filter = bson.M{"account": session.Account}
	default:
		return nil, errors.New("Unsupported Key.")
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "account", Value: session.Account},
			{Key: "access_token", Value: session.AccessToken},
			{Key: "refresh_token", Value: session.RefreshToken},
			{Key: "created_at", Value: session.CreatedAt},
			{Key: "updated_at", Value: session.UpdatedAt},
		}},
	}

	// Belgeyi bulup güncelle, eğer bulunamazsa ekle
	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Matched %v document(s) and modified %v document(s)\n", result.MatchedCount, result.ModifiedCount)
	return session, nil
}
