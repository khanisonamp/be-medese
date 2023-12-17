package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Order struct {
	AnOrderId             *uuid.UUID `json:"an_order_id"`
	USerId                *uuid.UUID `json:"user_id"`
	SpOrderParcelId       *uuid.UUID `json:"sp_order_parcel_id"`
	ReferenceNo           string     `json:"reference_no"`
	DesName               string     `json:"des_name"`
	DesPhoneNumber        string     `json:"des_phone_number"`
	DesAddress            string     `json:"des_address"`
	DesSubdistrict        string     `json:"des_subdistrict"`
	DesProvince           string     `json:"des_province"`
	DesPostcode           string     `json:"des_postcode"`
	CreatedAt             *time.Time `json:"created_at"`
	DeletedAt             int64      `json:"deleted_at"`
	CourierCode           int16      `json:"courier_code"`
	CodAmount             int        `json:"cod_amount"`
	FulfillmentStatus     int16      `json:"fulfillment_status"`
	ShippingStatus        string     `json:"shipping_status"`
	CodStatus             string     `json:"cod_status"`
	StatusCompletedDate   *time.Time `json:"status_completed_date"`
	CodTransferredDate    *time.Time `json:"cod_transferred_date"`
	PackagedImagegUrl     string     `json:"packaged_imageg_url"`
	TrackingCode          string     `json:"tracking_code"`
	JnaCodTransferredDate string     `json:"jna_cod_transferred_date"`
	SortCode              string     `json:"sort_code"`
	LineCode              string     `json:"line_code"`
	SortingLineCode       string     `json:"sorting_line_code"`
	DstStoreName          string     `json:"dst_store_name"`
	DesDistrict           string     `json:"des_district"`
}
