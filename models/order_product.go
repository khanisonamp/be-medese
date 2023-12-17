package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderProduct struct {
	AnOrderProductId *uuid.UUID `json:"an_order_product_id"`
	AnOrderId        *uuid.UUID `json:"an_order_id"`
	AnProductId      *uuid.UUID `json:"an_product_id"`
	Quantity         int        `json:"quantity"`
	CreatedAt        *time.Time `json:"created_at"`
	DeletedAt        int64      `json:"deleted_at"`
	Note             string     `json:"note"`
}
