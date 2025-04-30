package product

type Product struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Name2    string `json:"name2" bson:"name2"`
	Price    int    `json:"price" bson:"price"`
	Category string `json:"category" bson:"category"`
}
