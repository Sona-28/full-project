package models

type Inventory struct {
	ID          int64  `json:"id" bson:"_id"`
	Item 	  string `json:"item" bson:"item"`
	Features []string `json:"features" bson:"features"`
	Categories []string `json:"categories" bson:"categories"`
	Skus []Sku `json:"skus" bson:"skus"`
}