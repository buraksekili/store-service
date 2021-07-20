package api

import "github.com/buraksekili/store-service/products"

type addProductRes struct {
	ID string `json:"id"`
}

type productPage struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	ImageURL    string             `json:"image_url"`
	Price       float32            `json:"price"`
	Stock       int                `json:"stock"`
	VendorName  string             `json:"vendor_name"`
	Comments    []products.Comment `json:"comments,omitempty"`
}

type getProductRes struct {
	Product productPage `json:"product"`
	Error   string      `json:"error,omitempty"`
}

type listProductsRes struct {
	Products []productPage `json:"products"`
	Error    string        `json:"error,omitempty"`
}

type listVendorProductsRes struct {
	Products []productPage `json:"products"`
	Error    string        `json:"error,omitempty"`
}
