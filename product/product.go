package product

// struct
type Product struct {
	Id		 int	`json:"id"`
	Name     string	`json:"name"`
	Price    int	`json:"price"`
	Category string	`json:"category"`
}