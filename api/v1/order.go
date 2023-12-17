package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"
	"fmt"
	"strconv"
	"time"

	helper "api-medese/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func OrderGetAll(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	//request
	payload := struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid body",
			MessageTh: "body ไม่ถูกต้อง",
		})
	}

	if payload.StartDate == "" || payload.EndDate == "" {
		dateNow := time.Now().Format("2006-01-02")
		payload.StartDate = dateNow + " 00:00"
		payload.EndDate = dateNow + " 23:59"
	}

	start, _ := time.Parse("2006-01-02 15:04", payload.StartDate)
	end, _ := time.Parse("2006-01-02 15:04", payload.EndDate)

	startD := helper.ParseLocal(start)
	endD := helper.ParseLocal(end)

	fmt.Println("startD : ", startD)
	fmt.Println("endD : ", endD)

	//get order
	var getOrder []models.Order = make([]models.Order, 0)

	order := db.Table("order").Unscoped().Where("created_at BETWEEN ? AND ?", startD, endD)
	if err := order.Order("created_at desc").Find(&getOrder).Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get order error",
			MessageTh: "get order error",
			Error:     "not found",
		})
	}

	no := 1
	result := make([]map[string]interface{}, 0)
	for _, v := range getOrder {
		var getOrderProduct []models.OrderProduct = make([]models.OrderProduct, 0)
		orderProduct := db.Table("order_product").Unscoped().Where("an_order_id=?", v.AnOrderId.String())
		if err := orderProduct.Order("created_at desc").Find(&getOrderProduct).Error; err != nil {
			logrus.Info("error --> ", err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(Result{
				Status:    fiber.StatusNotFound,
				Message:   "get orderProduct error",
				MessageTh: "get orderProduct error",
				Error:     "not found",
			})
		}

		var ProductCode string
		var quantity int
		var quantityStr string

		if len(getOrderProduct) > 0 {
			for _, v := range getOrderProduct {
				fmt.Println("v.AnOrderProductId : ", v.AnOrderProductId)
				var getProduct models.Product
				product := db.Table("product").Unscoped().Where("an_product_id=?", v.AnProductId.String())
				if err := product.Order("created_at desc").Find(&getProduct).Error; err != nil {
					logrus.Info("error --> ", err.Error())
					return ctx.Status(fiber.StatusNotFound).JSON(Result{
						Status:    fiber.StatusNotFound,
						Message:   "get Product error",
						MessageTh: "get Product error",
						Error:     "not found",
					})
				}

				if ProductCode != "" {
					convert := ProductCode + ",\n " + getProduct.ProductCode
					ProductCode = convert
				} else {
					ProductCode = getProduct.ProductCode
				}

				if quantityStr != "" {
					convertQuantity := quantityStr + ",\n " + strconv.Itoa(v.Quantity)
					quantityStr = convertQuantity
				} else {
					quantityStr = strconv.Itoa(v.Quantity)
				}

				quantity += v.Quantity
			}
		}

		resultData := map[string]interface{}{
			"no":               no,
			"an_order_id":      v.AnOrderId,
			"reference_no":     v.ReferenceNo,
			"product_name":     ProductCode,
			"quantity_sum":     quantity,
			"quantity":         quantityStr,
			"tracking_code":    v.TrackingCode,
			"des_name":         v.DesName,
			"des_phone_number": v.DesPhoneNumber,
		}

		result = append(result, resultData)

		no += 1
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      result,
	})
}
