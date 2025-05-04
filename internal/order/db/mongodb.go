package db

import (
	"context"
	"errors"
	"fmt"
	"server/internal/order"
	"server/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (s *orderStorage) Create(ctx context.Context, order order.Order) (string, error) {
	res, err := s.collection.InsertOne(ctx, order)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert objectid to hex")
	}

	return oid.Hex(), nil
}

func (s *orderStorage) GetByID(ctx context.Context, id string) (order.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return order.Order{}, fmt.Errorf("failed to convert hex to objectid: %v", err)
	}

	var o order.Order
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&o)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return order.Order{}, fmt.Errorf("order not found")
		}
		return order.Order{}, fmt.Errorf("failed to find order: %v", err)
	}

	return o, nil
}

func (s *orderStorage) GetByUserID(ctx context.Context, userID string) ([]order.Order, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find orders: %v", err)
	}
	defer cursor.Close(ctx)

	var orders []order.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}

	return orders, nil
}

func (s *orderStorage) UpdateStatus(ctx context.Context, id, status string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid: %v", err)
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (s *orderStorage) Cancel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectid: %v", err)
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func NewOrderStorage(database *mongo.Database, logger *logging.Logger) order.Storage {
	return &orderStorage{
		collection: database.Collection("orders"),
		logger:     logger,
	}
}
