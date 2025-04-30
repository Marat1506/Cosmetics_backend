package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/internal/product"
	"server/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *ProductStorage) Create(ctx context.Context, product product.Product) (string, error) {
	d.logger.Debug("create product")
	result, err := d.collection.InsertOne(ctx, product)
	if err != nil {
		return "", fmt.Errorf("failed to create product: %v", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", fmt.Errorf("failed to convert objectid to hex")
}

func (d *ProductStorage) GetAll(ctx context.Context) ([]product.Product, error) {
	d.logger.Debug("get all products")

	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := d.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []product.Product
	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			d.logger.Errorf("Failed to decode product: %v", err)
			continue
		}
		products = append(products, p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	d.logger.Debugf("Retrieved %d products", len(products))
	return products, nil
}

func NewStorage(database *mongo.Database, collectionName string, logger *logging.Logger) *ProductStorage {
	return &ProductStorage{
		collection: database.Collection(collectionName),
		logger:     logger,
	}
}
