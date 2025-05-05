package order

type Order struct {
	ID         string   `json:"id" bson:"id"`                 // ID заказа
	Products   []string `json:"products" bson:"products"`     // Массив ID товаров
	Status     string   `json:"status" bson:"status"`         // Статус заказа
	CreatedAt  int64    `json:"createdAt" bson:"createdAt"`   // Временная метка создания
	TotalPrice int      `json:"totalPrice" bson:"totalPrice"` // Общая сумма заказа
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
