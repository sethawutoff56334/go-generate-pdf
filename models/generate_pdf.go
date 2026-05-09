package model

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PhoneNo   string `json:"phone_no"`
}

type PricingSummary struct {
	SubTotal float64 `json:"sub_total"`
	Discount float64 `json:"discount"`
	Total    float64 `json:"total"`
}

type Product struct {
	ProductName  string  `json:"product_name"`
	Qty          int     `json:"qty"`
	PricePerUnit float64 `json:"price_per_unit"`
}

type GeneratePDFRequest struct {
	ReceiptNo      string         `json:"receipt_no"`
	Customer       Customer       `json:"customer"`
	PricingSummary PricingSummary `json:"pricing_summary"`
	ProductList    []Product      `json:"product_list"`
}

// Response structs

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	File    string `json:"file"`
}
