package models

type Customer struct{
	ID string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	Email string `json:"email" bson:"email"`
	Address AddressStruct `json:"address" bson:"address"`
	ShippingAddress AddressStruct `json:"shipping_address" bson:"shipping_address"`
}

type AddressStruct struct{
	Country string `json:"country" bson:"country"`
	Street1 string `json:"street1" bson:"street1"`
	Street2 string `json:"street2" bson:"street2"`
	City string `json:"city" bson:"city"`
	State string `json:"state" bson:"state"`
	Zip string `json:"zip" bson:"zip"`
}