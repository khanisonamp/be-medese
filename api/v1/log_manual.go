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

func GetLogManualOrder(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	var getLogManualOrder []models.LogManualOrder = make([]models.LogManualOrder, 0)
	if err := db.Model(&models.LogManualOrder{}).Find(&getLogManualOrder).Error; err != nil {
		logrus.Errorln("error getLogManualOrder : ", err)
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get Product error",
			MessageTh: "get Product error",
			Error:     err,
		})
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range getLogManualOrder {

		resultData := map[string]interface{}{
			"created_at":   v.CreatedAt.Format("2006-01-02 15:04"),
			"product_code": v.ProductCode,
			"product_name": v.ProductName,
			"order_amount": v.OrderAmount,
			"remark":       v.Remark,
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
