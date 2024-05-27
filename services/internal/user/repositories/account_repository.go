package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"services/internal/user/models"
)

// BookRepository MongoDB işlemleri için bir yapı
type AccountRepository struct {
	collection       *mongo.Collection
	findAllAggregate []bson.M
	indexModel       []mongo.IndexModel
}

type userObjectID struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

// NewBookRepository yeni bir BookRepository oluşturur
func NewAccountRepository(db *mongo.Database, collectionName string) (*AccountRepository, error) {

	

	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(collectionName) // Koleksiyon adınıza uygun olarak değiştirin

	var isUnique = true

	userUnique := mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: &options.IndexOptions{Unique: &isUnique}}

	indexModel := []mongo.IndexModel{userUnique}

	_, err = collection.Indexes().CreateMany(context.TODO(), indexModel, options.CreateIndexes())
	if err != nil {
		return nil, err
	}

	return &AccountRepository{
		collection: collection,
		findAllAggregate: []primitive.M{
			bson.M{"$lookup": bson.M{"from": "users", "localField": "user._id", "foreignField": "_id", "as": "user"}},
			bson.M{"$project": bson.M{"username": 1, "password": 1, "created_at": 1, "updated_at": 1, "user": bson.M{"$first": "$user"}}}},
	}, nil
}

// Create yeni bir kitap ekler
func (r *AccountRepository) Create(account *models.Account) (*models.Account, error) {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	res, err := r.collection.InsertOne(context.TODO(), account, options.InsertOne())
	if err != nil {
		return nil, err
	}
	account.ID = res.InsertedID.(primitive.ObjectID)

	return account, nil
}

// GetAll tüm kitapları getirir
func (r *AccountRepository) GetAll() ([]*models.Account, error) {

	//opt := options.Find()
	//cur, err := r.collection.Find(context.TODO(), nil, opt)
	cur, err := r.collection.Aggregate(context.TODO(), r.findAllAggregate, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var books []*models.Account
	for cur.Next(context.Background()) {
		var book models.Account
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (r *AccountRepository) Delete(accountID primitive.ObjectID) error {
	filter := bson.M{"_id": accountID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("account Not Found.")
	}
	return nil
}

// FindOne belirli bir hesabı ID ile bulur
func (r *AccountRepository) FindOne(accountID primitive.ObjectID) (*models.Account, error) {
	filter := bson.M{"_id": accountID}
	var account models.Account
	err := r.collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *AccountRepository) FindOneWithParameters(key string, data any) (*models.Account, error) {
	filter := bson.M{key: data}
	var account models.Account
	err := r.collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindAll tüm hesapları getirir
func (r *AccountRepository) FindAll() ([]*models.Account, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var accounts []*models.Account
	for cur.Next(context.Background()) {
		var account models.Account
		err := cur.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// Update bir hesabın bilgilerini günceller
func (r *AccountRepository) Update(accountID primitive.ObjectID, updatedAccount *models.Account) error {
	filter := bson.M{"_id": accountID}
	updatedAccount.UpdatedAt = time.Now()
	update := bson.M{
		"$set": updatedAccount,
	}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Account Not Found.")
	}
	return nil
}
