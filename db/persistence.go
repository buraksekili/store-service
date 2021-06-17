package db

type DBHandler interface {
	AddProduct(Product) ([]byte, error)
	AddVendor(Vendor) ([]byte, error)
	AddComment(Comment, string) ([]byte, error)
	AddUser(User) ([]byte, error)

	FindProduct(string) (Product, error)
	FindProductsByVendor(string) ([]Product, error)
	FindUser(string) (User, error)
	FindUserByName(string) (User, error)

	GetProducts() ([]Product, error)
	GetComments() ([]Comment, error)
	GetVendors() ([]Vendor, error)
	GetUsers() ([]User, error)
}
