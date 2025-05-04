package user

type User struct {
	ID           string   `json:"id" bson:"_id,omitempty"`
	Email        string   `json:"email" bson:"email"`
	Username     string   `json:"username" bson:"username"`
	PasswordHash string   `json:"-" bson:"password"`
	Favorites    []string `json:"favorites" bson:"favorites"`
	Cart         []string `json:"cart" bson:"cart"`
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateFavoritesDTO struct {
	ProductID string `json:"productId"`
}

type UpdateCartDTO struct {
	ProductID string `json:"productId"`
}
