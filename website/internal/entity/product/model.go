package product

type Product struct {
	Id                int
	Name              string
	FullName          string
	Waight            int //gramm
	Category          ProductCategory
	UnitOfMeasurement string
	Remains           int
	Price             int
	Discount          DiscountProduct
}

type ProductCategory struct {
	Name string
}

type DiscountProduct struct {
	Percent           int
	PriceWithDiscount int
	Discount          bool
}
