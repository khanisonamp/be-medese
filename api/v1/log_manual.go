package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CreatedLogManual(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	//request
	payload := struct {
		ProductCode string `json:"product_code"`
		OrderAmount string `json:"order_amount"`
		Remark      string `json:"remark"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid body",
			MessageTh: "body ไม่ถูกต้อง",
		})
	}

	var getProduct models.Product
	product := db.Table("product").Unscoped().Where("product_code=?", payload.ProductCode)
	if err := product.Order("created_at desc").Find(&getProduct).Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get Product error",
			MessageTh: "get Product error",
			Error:     "not found",
		})
	}

	setCreatedLogManual := models.LogManualOrder{
		ProductCode: getProduct.ProductCode,
		ProductName: getProduct.Name,
		OrderAmount: payload.OrderAmount,
		Remark:      payload.Remark,
	}

	db.Model(&models.LogManualOrder{}).Create(&setCreatedLogManual)

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
	})
}
