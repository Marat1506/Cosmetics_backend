package order

type Order struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"username"`
	Phone        string `json:"phone" bson:"phone"`
	TelegramNick string `json:"telegram_nick" bson:"telegram_nick"`
	Completed    bool   `json:"completed" bson:"completed"`
}

type CreateOrderDTO struct {
	Username     string `json:"username" bson:"username"`
	Phone        string `json:"phone" bson:"phone"`
	TelegramNick string `json:"telegram_nick" bson:"telegram_nick"`
}
