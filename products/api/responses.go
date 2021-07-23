package api

import "github.com/buraksekili/store-service/products"

type addProductRes struct {
	ID string `json:"id"`
}

// type productPage struct {
// 	ID          string             `json:"id"`
// 	Name        string             `json:"name"`
// 	Category    string             `json:"category"`
// 	Description string             `json:"description"`
// 	ImageURL    string             `json:"image_url"`
// 	Price       float32            `json:"price"`
// 	Stock       int                `json:"stock"`
// 	VendorName  string             `json:"vendor_name"`
// 	Comments    []products.Comment `json:"comments,omitempty"`
// }

type getProductRes struct {
	Product products.Product `json:"product"`
}

type listProductsRes struct {
	Total    int                `json:"total"`
	Offset   int                `json:"offset"`
	Limit    int                `json:"limit"`
	Products []products.Product `json:"products"`
}

type listVendorProductsRes struct {
	Products []products.Product `json:"products"`
}

type getAllCommentsRes struct {
	Total    int                `json:"total"`
	Offset   int                `json:"offset"`
	Limit    int                `json:"limit"`
	Comments []products.Comment `json:"comments"`
}

type getProductCommentsRes struct {
	Comments []products.Comment `json:"comments"`
}
