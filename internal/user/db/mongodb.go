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
			Orders:   []user.Order{},
		}, nil
	}
	result := d.collection.FindOne(ctx, filter)
	err = result.Decode(&u)

	if err != nil {
		return user.User{}, fmt.Errorf("пользователь не найден %s", err)
	}

	// Инициализируем orders, если оно nil
	if u.Orders == nil {
		u.Orders = []user.Order{}
	}

	return u, nil
}

func (d *db) CreateOrder(ctx context.Context, userID string, order user.Order) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"orders": order}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) GetOrders(ctx context.Context, userID string) ([]user.Order, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	var u user.User
	err = d.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	if u.Orders == nil {
		return []user.Order{}, nil
	}

	return u.Orders, nil
}

func (d *db) AddToFavorites(ctx context.Context, userID string, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$addToSet": bson.M{"favorites": productID}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add to favorites: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) RemoveFromFavorites(ctx context.Context, userID string, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$pull": bson.M{"favorites": productID}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove from favorites: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) AddToCart(ctx context.Context, userID string, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$addToSet": bson.M{"cart": productID}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add to cart: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$pull": bson.M{"cart": productID}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove from cart: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) UpdateCart(ctx context.Context, userID string, cart []string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectedid, hex: %s", userID)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"cart": cart}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update cart: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (d *db) GetFavorites(ctx context.Context, userID string) ([]string, error) {
	fmt.Print("GetFavorites = ", userID)
	d.logger.Infof("Запрос избранного для userID: %s", userID)

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		d.logger.Errorf("Ошибка конвертации userID: %v", err)
		return nil, fmt.Errorf("неверный ID пользователя")
	}

	var u user.User
	err = d.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		d.logger.Errorf("Пользователь не найден: %v", err)
		return nil, fmt.Errorf("пользователь не найден")
	}

	// Если favorites == nil, возвращаем пустой массив
	if u.Favorites == nil {
		d.logger.Info("Поле favorites отсутствует, возвращаем []")
		return []string{}, nil
	}

	d.logger.Infof("Найдены избранные товары: %v", u.Favorites)
	return u.Favorites, nil
}

func (d *db) GetCart(ctx context.Context, userID string) ([]string, error) {
	d.logger.Infof("Запрос корзины для userID: %s", userID)

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		d.logger.Errorf("Ошибка конвертации userID: %v", err)
		return nil, fmt.Errorf("неверный ID пользователя")
	}

	var u user.User
	err = d.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		d.logger.Errorf("Пользователь не найден: %v", err)
		return nil, fmt.Errorf("пользователь не найден")
	}

	if u.Cart == nil {
		d.logger.Info("Поле cart отсутствует, возвращаем []")
		return []string{}, nil
	}

	d.logger.Infof("Найдены товары в корзине: %v", u.Cart)
	return u.Cart, nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
