package products

// Product -
type Product struct {
	ID            string  `json:"id" docstore:"id"`
	Sku           string  `json:"sku" docstore:"sku" validate:"required"`
	Title         string  `json:"title" docstore:"title" validate:"required,gte=1,lte=100"`
	Condition     string  `json:"condition" docstore:"condition" validate:"required,gte=1,lte=50"`
	NumberInStock uint32  `json:"numberInStock" docstore:"numberInStock" validate:"required,gte=0,lte=10000"`
	UnitCost      float64 `json:"unitCost" docstore:"unitCost" validate:"required,gte=1,lte=100000"`
}

// UpdateProduct -
type UpdateProduct struct {
	ID            string  `json:"id" docstore:"id"`
	Sku           string  `json:"sku" docstore:"sku" validate:"required"`
	Title         string  `json:"title" docstore:"title" validate:"required,gte=1,lte=100"`
	Condition     string  `json:"condition" docstore:"condition" validate:"required,gte=1,lte=50"`
	NumberInStock uint32  `json:"numberInStock" docstore:"numberInStock" validate:"required,gte=0,lte=10000"`
	UnitCost      float64 `json:"unitCost" docstore:"unitCost" validate:"required,gte=1,lte=100000"`
}
