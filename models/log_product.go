package models

import uuid "github.com/satori/go.uuid"

type LogProduct struct {
	Model
	AnProductId    *uuid.UUID `json:"an_product_id"`
	ProductCode    string     `json:"product_code"`
	RemainingStock string     `json:"remaining_stock"`
	OrderInSystem  string     `json:"order_in_system"`
	OrderOutSystem string     `json:"order_out_system"`
	RemainingToday string     `json:"remaining_today"`
	DateTxt        string     `json:"date_txt"`
}

type LogStock struct {
	Model
	ProductCode    string `json:"product_code"`
	ProductName    string `json:"product_name"`
	RemainingStock string `json:"remaining_stock"`
	StockIn        string `json:"stock_in"`
	RemainingToday string `json:"remaining_today"`
}

type LogManualOrder struct {
	Model
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	OrderAmount string `json:"order_amount"`
	Remark      string `json:"remark"`
}
