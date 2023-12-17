package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	AnProductId *uuid.UUID `json:"an_product_id"`
	UserId      *uuid.UUID `json:"user_id"`
	ProductCode string     `json:"product_code"`
	Name        string     `json:"name"`
	CreatedAt   *time.Time `json:"created_at"`
	DeletedAt   int64      `json:"deleted_at"`
	ImgUrl      string     `json:"img_url"`
	SerialRegex string     `json:"serial_regex"`
	RoboticSku  string     `json:"robotic_sku"`
	Stock       string     `json:"stock"`
}
