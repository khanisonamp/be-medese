package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"
	"fmt"
	"time"

	helper "api-medese/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetLogProduct(ctx *fiber.Ctx) error {
	db := dbCon.DBConn
	_ = db

	//request
	month := ctx.Query("month")

	if month == "" {
		dateNow := time.Now().Format("2006-01")
		month = dateNow
	}

	monthStr := month + "-01 00:00"
	monthEnd := month + "-31 23:59"

	start, _ := time.Parse("2006-01-02 15:04", monthStr)
	end, err := time.Parse("2006-01-02 15:04", monthEnd)

	if err != nil {
		logrus.Info("charge day", err)
		monthEnd := month + "-30 23:59"
		end, _ = time.Parse("2006-01-02 15:04", monthEnd)
	}

	startD := helper.ParseLocal(start)
	endD := helper.ParseLocal(end)

	fmt.Println("startD : ", startD)
	fmt.Println("endD : ", endD)

	var getLogProduct []models.LogProduct = make([]models.LogProduct, 0)
	if err := db.Model(&models.LogProduct{}).Where("created_at BETWEEN ? AND ?", startD, endD).Order("created_at desc").Find(&getLogProduct).Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get getLogProduct error",
			MessageTh: "get getLogProduct error",
			Error:     "not found",
		})
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range getLogProduct {

		resultData := map[string]interface{}{
			"product_code":     v.ProductCode,
			"remaining_stock":  v.RemainingStock,
			"order_in_system":  v.OrderInSystem,
			"order_out_system": v.OrderOutSystem,
			"remaining_today":  v.RemainingToday,
			"date_txt":         v.DateTxt,
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
