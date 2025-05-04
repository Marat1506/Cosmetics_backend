package order

import "time"

type Order struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	UserID     string    `json:"userId" bson:"userId"`
	Products   []Product `json:"products" bson:"products"`
	Status     string    `json:"status" bson:"status"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalPrice int       `json:"totalPrice" bson:"totalPrice"`
}

type Product struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Price    int    `json:"price" bson:"price"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

type CreateOrderDTO struct {
	UserID     string    `json:"userId"`
	Products   []Product `json:"products"`
	TotalPrice int       `json:"totalPrice"`
}

type UpdateStatusDTO struct {
	Status string `json:"status"`
}
