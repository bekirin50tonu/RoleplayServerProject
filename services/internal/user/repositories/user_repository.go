package repositories

import (
	"context"
	"errors"
	"services/internal/user/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) (*UserRepository, error) {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}

	collection := db.Collection(collectionName)
	var isUnique = true

	userUnique := mongo.IndexModel{Keys: bson.D{{Key: "email", Value: 1}}, Options: &options.IndexOptions{Unique: &isUnique}}

	indexModel := []mongo.IndexModel{userUnique}

	_, err = collection.Indexes().CreateMany(context.TODO(), indexModel, options.CreateIndexes())
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		collection: *collection,
	}, nil
}

func (r *UserRepository) Create(User *models.User) (*models.User, error) {
	User.CreatedAt = time.Now()
	User.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(context.TODO(), User)
	if err != nil {
		return nil, err
	}
	User.ID = res.InsertedID.(primitive.ObjectID)

	return User, nil
}

// GetAll tüm kitapları getirir
func (r *UserRepository) GetAll() ([]*models.User, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var books []*models.User
	for cur.Next(context.Background()) {
		var book models.User
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (r *UserRepository) Delete(UserID primitive.ObjectID) error {
	filter := bson.M{"_id": UserID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("User Not Found.")
	}
	return nil
}

// FindOne belirli bir hesabı ID ile bulur
func (r *UserRepository) FindOne(UserID primitive.ObjectID) (*models.User, error) {
	filter := bson.M{"_id": UserID}
	var User models.User
	err := r.collection.FindOne(context.TODO(), filter).Decode(&User)
	if err != nil {
		return nil, err
	}
	return &User, nil
}
func (r *UserRepository) FindOneWithParameter(key string, data any) (*models.User, error) {
	filter := bson.M{key: data}
	var User models.User
	err := r.collection.FindOne(context.TODO(), filter).Decode(&User)
	if err != nil {
		return nil, err
	}
	return &User, nil
}

// FindAll tüm hesapları getirir
func (r *UserRepository) FindAll() ([]*models.User, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var Users []*models.User
	for cur.Next(context.Background()) {
		var User models.User
		err := cur.Decode(&User)
		if err != nil {
			return nil, err
		}
		Users = append(Users, &User)
	}

	return Users, nil
}

// Update bir hesabın bilgilerini günceller
func (r *UserRepository) Update(UserID primitive.ObjectID, updatedUser *models.User) error {
	filter := bson.M{"_id": UserID}
	updatedUser.UpdatedAt = time.Now()
	update := bson.M{
		"$set": updatedUser,
	}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("User Not Found.")
	}
	return nil
}
