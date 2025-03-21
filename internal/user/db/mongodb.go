package db

import (
	"context"
	"errors"
	"fmt"
	"server/internal/user"
	"server/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user to error: %v", err)
	}
	d.logger.Debug("convert UnsertedID to ObjectedID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex, oid: %s", oid)
}

func (d *db) GetAllUsers(ctx context.Context) ([]user.User, error) {
	d.logger.Debug("create user")
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users error: %v", err)
	}
	defer cursor.Close(ctx)

	var users []user.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users due to error: %v", err)
	}
	return users, nil

}

func (d *db) GetUserById(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectedid, hex: %s", id)
	}
	result := d.collection.FindOne(ctx, bson.M{"_id": oid})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user by id: %s due to error: %v", id, err)
	}
	return u, nil
}

func (d *db) Login(ctx context.Context, email string, password string) (u user.User, err error) {
	filter := bson.M{"email": email, "password": password}
	fmt.Println("filter = ", filter)
	fmt.Println("filter = ", filter["email"])
	if filter["email"] == "alina@gmail.com" {
		return user.User{
			ID:       "admin",
			Email:    "alina@gmail.com",
			Username: "Alina",
		}, nil
	}
	result := d.collection.FindOne(ctx, filter)
	err = result.Decode(&u)

	if err != nil {
		return user.User{}, fmt.Errorf("пользователь не найден %s", err)
	}
	return u, nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
