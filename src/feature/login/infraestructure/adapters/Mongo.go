package adapters

import (
	"context"
	"time"

	"ApiMedicGO/src/feature/login/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserRepository implementa repositories.UserRepository usando MongoDB.
type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) FindByLicenseNumber(licenseNumber string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"license_number": licenseNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser crea un nuevo usuario en la base de datos.
func (r *MongoUserRepository) CreateUser(user *entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}
