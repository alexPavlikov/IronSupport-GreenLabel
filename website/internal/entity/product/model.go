package product

type Product struct {
	Id                int
	Article           string
	Name              string
	FullName          string
	Waight            int //gramm
	Category          ProductCategory
	UnitOfMeasurement string
	Remains           int
	Price             int
	Discount          DiscountProduct
	OnTheWay          bool
}

type ProductCategory struct {
	Name        string
	Avatar      string
	Description string
	MinAge      int
}

type DiscountProduct struct {
	Percent           int
	PriceWithDiscount int
	Discount          bool
}
