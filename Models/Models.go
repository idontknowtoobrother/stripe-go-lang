package models

type Product struct {
	Uuid        string  `bson:"uuid" json:"uuid"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Price       float64 `bson:"price" json:"price"`
}

func (b *Product) GetCollectionName() string {
	return "products"
}
