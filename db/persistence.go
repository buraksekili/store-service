package db

type DBHandler interface {
	AddProduct(Product) ([]byte, error)
	// TODO: implement
	// AddVendor(Vendor) ([]byte, error)
	// AddComment(Comment, Product) ([]byte, error)
	//
	// FindProduct(string) (Product, error)
	// FindProductsByVendor(string) ([]Product, error)
	// FindVendor(string) (Vendor, error)
	//
	GetProducts() ([]Product, error)
}
