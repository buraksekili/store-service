package db

type DBHandler interface {
	AddProduct(Product) ([]byte, error)
	AddVendor(Vendor) ([]byte, error)
	AddComment(Comment, string) ([]byte, error)

	FindProduct(string) (Product, error)
	FindProductsByVendor(string) ([]Product, error)

	GetProducts() ([]Product, error)
	GetComments() ([]Comment, error)
	GetVendors() ([]Vendor, error)
}
