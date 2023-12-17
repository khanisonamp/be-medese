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

type OrderProduct struct {
	No          string `json:"No"`
	AnProductId string `json:"an_product_id"`
	Quantity    int    `json:"quantity"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
}

func OrderProductDashBoard(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

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

	//get order product
	var getOrderProduct []models.OrderProduct = make([]models.OrderProduct, 0)
	orderProduct := db.Table("order_product").Unscoped().Where("created_at BETWEEN ? AND ?", startD, endD)
	if err := orderProduct.Order("created_at desc").Find(&getOrderProduct).Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get OrderProduct error",
			MessageTh: "get OrderProduct error",
			Error:     "not found",
		})
	}

	unique := make(map[string]bool, 0)
	newOrderProduct := make([]OrderProduct, 0)
	var i int
	for _, v := range getOrderProduct {
		if !unique[v.AnProductId.String()] {
			unique[v.AnProductId.String()] = true
			i += 1
			setOrderProduct := OrderProduct{}
			setOrderProduct.No = strconv.Itoa(i)
			setOrderProduct.AnProductId = v.AnProductId.String()
			setOrderProduct.ProductName = v.Note
			setOrderProduct.Quantity = 0
			newOrderProduct = append(newOrderProduct, setOrderProduct)
		}

	}

	for i, x := range newOrderProduct {

		quantityNew := x.Quantity
		for _, v := range getOrderProduct {
			quanStr := strconv.Itoa(v.Quantity)
			quanInt, _ := strconv.Atoi(quanStr)

			if x.AnProductId == v.AnProductId.String() {
				quantityNew += quanInt
				newOrderProduct[i].Quantity = quantityNew
				fmt.Println("s===> ", quantityNew)
			}
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range newOrderProduct {

		var getProduct models.Product
		product := db.Table("product").Unscoped().Where("an_product_id=?", v.AnProductId)
		if err := product.Order("created_at desc").Find(&getProduct).Error; err != nil {
			logrus.Info("error --> ", err.Error())
			return ctx.Status(fiber.StatusNotFound).JSON(Result{
				Status:    fiber.StatusNotFound,
				Message:   "get Product error",
				MessageTh: "get Product error",
				Error:     "not found",
			})
		}

		resultData := map[string]interface{}{
			"no":            v.No,
			"an_product_id": v.AnProductId,
			"quantity":      v.Quantity,
			"product_code":  getProduct.ProductCode,
			"product_name":  getProduct.Name,
			"stock_total":   getProduct.Stock,
		}

		result = append(result, resultData)
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      result,
	})
}
