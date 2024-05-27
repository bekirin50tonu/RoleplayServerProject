package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"services/internal/storage/models"
)

// BookRepository MongoDB işlemleri için bir yapı
type StorageRepository struct {
	collection *mongo.Collection
}

// NewBookRepository yeni bir BookRepository oluşturur
func NewStorageRepository(db *mongo.Database, collectionName string) (*StorageRepository, error) {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		cmdErr, ok := err.(mongo.CommandError)
		if !ok {
			return nil, err
		}
		log.Fatalln("Create Collection:", cmdErr.Message)
	}
	collection := db.Collection(collectionName) // Koleksiyon adınıza uygun olarak değiştirin
	return &StorageRepository{
		collection: collection,
	}, nil
}

// Create yeni bir kitap ekler
func (r *StorageRepository) Create(Storage *models.Storage) (*models.Storage, error) {
	Storage.CreatedAt = time.Now()
	Storage.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(context.TODO(), Storage)
	if err != nil {
		return nil, err
	}
	Storage.ID = res.InsertedID.(primitive.ObjectID)

	return Storage, nil
}

// GetAll tüm kitapları getirir
func (r *StorageRepository) GetAll() ([]*models.Storage, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var books []*models.Storage
	for cur.Next(context.Background()) {
		var book models.Storage
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (r *StorageRepository) Delete(StorageID primitive.ObjectID) error {
	filter := bson.M{"_id": StorageID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Storage Not Found.")
	}
	return nil
}

// FindOne belirli bir hesabı ID ile bulur
func (r *StorageRepository) FindOne(StorageID primitive.ObjectID) (*models.Storage, error) {
	filter := bson.M{"_id": StorageID}
	var Storage models.Storage
	err := r.collection.FindOne(context.TODO(), filter).Decode(&Storage)
	if err != nil {
		return nil, err
	}
	return &Storage, nil
}
func (r *StorageRepository) FindOneWithParameters(filter bson.M) (*models.Storage, error) {
	var Storage models.Storage
	err := r.collection.FindOne(context.TODO(), filter).Decode(&Storage)
	if err != nil {
		return nil, err
	}
	return &Storage, nil
}

// FindAll tüm hesapları getirir
func (r *StorageRepository) FindAll() ([]*models.Storage, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var Storages []*models.Storage
	for cur.Next(context.Background()) {
		var Storage models.Storage
		err := cur.Decode(&Storage)
		if err != nil {
			return nil, err
		}
		Storages = append(Storages, &Storage)
	}

	return Storages, nil
}

// Update bir hesabın bilgilerini günceller
func (r *StorageRepository) Update(StorageID primitive.ObjectID, updatedStorage *models.Storage) error {
	filter := bson.M{"_id": StorageID}
	updatedStorage.UpdatedAt = time.Now()
	update := bson.M{
		"$set": updatedStorage,
	}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Storage Not Found.")
	}
	return nil
}
