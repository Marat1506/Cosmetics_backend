package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/order"
	"server/pkg/logging"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *OrderStorage) Create(ctx context.Context, order order.Order) (string, error) {
	d.logger.Debug("create order")
	result, err := d.collection.InsertOne(ctx, order)

	if err != nil {
		return "", fmt.Errorf("failed to create order to error: %v", err)
	}
	d.logger.Debug("convert UnsertedID to ObjectedID")
	oid, ok := result.InsertedID.(primitive.ObjectID)

	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(order)
	return "", fmt.Errorf("failed to convert objectid to hex, oid: %s", oid)
}

func (d *OrderStorage) GetOrders(ctx context.Context) ([]order.Order, error) {
	d.logger.Debug("get orders")
	cursor, err := d.collection.Find(ctx, bson.M{})
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

func (d *OrderStorage) ChangeOrder(ctx context.Context, id string) (o order.Order, err error) {
	d.logger.Debug("change order")
	fmt.Println("ID = ", id)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return order.Order{}, err
	}
	result := d.collection.FindOne(ctx, bson.M{"_id": oid})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return o, fmt.Errorf("not found")
		}
		return o, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)

	}
	if err := result.Decode(&o); err != nil {
		return o, fmt.Errorf("failed to decode user by id: %s due to error: %v", id, err)
	}
	fmt.Println("order", o.Completed)

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	newCopleted := !o.Completed
	err = d.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"completed": newCopleted}},
		opts).Decode(&o)

	if err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return o, fmt.Errorf("order not found")
		}
		return o, fmt.Errorf("failed to update order: %v", err)
	}

	fmt.Println("update = ", o)
	return o, nil
}

func (d *OrderStorage) DeleteOrder(ctx context.Context, id string) error {
	d.logger.Debug("delete order")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result := d.collection.FindOneAndDelete(ctx, bson.M{"_id": oid})
	if result.Err() != nil {
		return fmt.Errorf("failed to delete order: %v", result.Err())
	}
	return nil
}

func NewStorage(database *mongo.Database, collectionName string, logger *logging.Logger) *OrderStorage {
	fmt.Println("lllLLL = ", collectionName)
	return &OrderStorage{
		collection: database.Collection(collectionName),
		logger:     logger,
	}
}
