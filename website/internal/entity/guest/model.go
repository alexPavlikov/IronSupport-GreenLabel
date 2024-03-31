package guest

import "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"

type Guests struct { //table
	Id           int
	Email        string
	Firstname    string
	Lastname     string
	Patronymic   string
	Phone        string
	Password     string
	Age          uint8
	Organization Organization
	SaveCard     Card
	AllPurchase  []Purchase
	Banned       bool
}

type Card struct {
	Number string
	Date   string
	CVV    string
	Bank   string
}

type Purchase struct { //table
	Guests  string
	Product []product.Product
	Cost    uint
	Date    string
}

type Organization struct {
	INN      int
	KPP      int
	OGRN     int
	Name     string
	Fullname string
	City     string
	Country  string
}
