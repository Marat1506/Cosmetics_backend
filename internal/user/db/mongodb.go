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

type userStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (s *userStorage) Create(ctx context.Context, user user.User) (string, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert objectid to hex")
	}

	return oid.Hex(), nil
}

func (s *userStorage) GetByID(ctx context.Context, id string) (user.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to convert hex to objectid: %v", err)
	}

	var u user.User
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, fmt.Errorf("user not found")
		}
		return user.User{}, fmt.Errorf("failed to find user: %v", err)
	}

	return u, nil
}

func (s *userStorage) GetByEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, fmt.Errorf("user not found")
		}
		return user.User{}, fmt.Errorf("failed to find user: %v", err)
	}

	return u, nil
}

func (s *userStorage) AddToFavorites(ctx context.Context, userID, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert user id: %v", err)
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$addToSet": bson.M{"favorites": productID}},
	)
	return err
}

func (s *userStorage) RemoveFromFavorites(ctx context.Context, userID, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert user id: %v", err)
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$pull": bson.M{"favorites": productID}},
	)
	return err
}

func (s *userStorage) AddToCart(ctx context.Context, userID, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert user id: %v", err)
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$addToSet": bson.M{"cart": productID}},
	)
	return err
}

func (s *userStorage) RemoveFromCart(ctx context.Context, userID, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert user id: %v", err)
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$pull": bson.M{"cart": productID}},
	)
	return err
}

func NewUserStorage(database *mongo.Database, logger *logging.Logger) user.Storage {
	return &userStorage{
		collection: database.Collection("users"),
		logger:     logger,
	}
}
