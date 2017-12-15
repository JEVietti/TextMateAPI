package models

import (
	"time"
)

type Product struct {
	ID      int16     `json:"id"`
	Name    string    `json:"name"`
	Link    string    `json:"link"`
	Price   float32   `json:"price"`
	Posted  time.Time `json:"posted"`
	Barcode string    `json:"barcode"`
}

type Products []Product
