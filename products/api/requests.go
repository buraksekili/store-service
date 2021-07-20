package api

type addProductReq struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	VendorID    string  `json:"vendor_id"`
}

type getProductReq struct {
	ProductID string `json:"product_id"`
}

type listProductsReq struct {
	offset int
	limit  int
}

type listVendorProductsReq struct {
	offset   int
	limit    int
	VendorID string `json:"vendor_id"`
}
