package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDBURL string

	// Формируем URL для подключения
	if strings.HasPrefix(host, "mongodb+srv://") {
		// Если хост уже содержит SRV-запись, используем его как есть
		mongoDBURL = fmt.Sprintf("%s/%s?retryWrites=true&w=majority&appName=Cluster0&authSource=%s", host, database, authDB)
	} else {
		// Для локального MongoDB
		if username == "" && password == "" {
			mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
		} else {
			mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
		}
	}

	fmt.Println("MongoDB URL:", mongoDBURL)

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb due to error: %v", err)
	}

	// Проверка подключения
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb due to error: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client.Database(database), nil
}
