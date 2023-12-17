package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"
	"fmt"
	"strconv"
	"time"

	helper "api-medese/helper"

	notifyLine "api-medese/helper/notify_line"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func ScriptCrestedLogStockAuto(ctx *fiber.Ctx) error {

	CrestedLogStockAuto()

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
	})
}

func CrestedLogStockAuto() error {
	db := dbCon.DBConn

	logrus.Info("Start Stock Cron Job", time.Now())

	dateStartNow := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dateStart := dateStartNow + " 12:00"
	dateEndNow := time.Now().Format("2006-01-02")
	dateEnd := dateEndNow + " 12:00"

	start, _ := time.Parse("2006-01-02 15:04", dateStart)
	end, _ := time.Parse("2006-01-02 15:04", dateEnd)

	startD := helper.ParseLocal(start)
	endD := helper.ParseLocal(end)

	//get order product
	var getOrderProduct []models.OrderProduct = make([]models.OrderProduct, 0)
	orderProduct := db.Table("order_product").Unscoped().Where("created_at BETWEEN ? AND ?", startD, endD)
	if err := orderProduct.Order("created_at desc").Find(&getOrderProduct).Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return nil
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
			}
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range newOrderProduct {
		dateNowTxt := time.Now().Format("2006-01-02")

		var getProduct models.Product
		product := db.Table("product").Unscoped().Where("an_product_id=?", v.AnProductId)
		if err := product.Order("created_at desc").Find(&getProduct).Error; err != nil {
			logrus.Info("get product error --> ", err.Error())
			return nil
		}

		//chack create log product
		var getLogProduct []models.LogProduct = make([]models.LogProduct, 0)
		logProduct := db.Model(&models.LogProduct{}).Where("date_txt=? AND an_product_id=?", dateNowTxt, v.AnProductId).First(&getLogProduct)
		_ = logProduct
		if len(getLogProduct) != 0 {
			fmt.Println("continue !!!")
			continue
		}

		//get order manual
		var getLogManualOrder []models.LogManualOrder = make([]models.LogManualOrder, 0)
		if err := db.Model(&models.LogManualOrder{}).Where("product_code = ? AND created_at BETWEEN ? AND ?", getProduct.ProductCode, startD, endD).Find(&getLogManualOrder).Error; err != nil {
			logrus.Errorln("get LogManualOrder error --> ", err)
		}

		orderManualAmount := 0
		if len(getLogManualOrder) > 0 {
			for _, v := range getLogManualOrder {
				orderAmount, _ := strconv.Atoi(v.OrderAmount)

				orderManualAmount += orderAmount
			}
		}
		orderManualAmountStr := strconv.Itoa(orderManualAmount)
		sumOrder := orderManualAmount + v.Quantity

		anProductIdUid := uuid.FromStringOrNil(v.AnProductId)
		quantityStr := strconv.Itoa(v.Quantity)

		//sum stock
		stockTodayInt, _ := strconv.Atoi(getProduct.Stock)
		remainingToday := stockTodayInt - sumOrder

		remainingTodayStr := strconv.Itoa(remainingToday)

		//created log product
		setCreatedLogLogProduct := models.LogProduct{
			AnProductId:    &anProductIdUid,
			ProductCode:    getProduct.ProductCode,
			RemainingStock: getProduct.Stock,
			OrderInSystem:  quantityStr,
			OrderOutSystem: orderManualAmountStr,
			RemainingToday: remainingTodayStr,
			DateTxt:        dateNowTxt,
		}

		createLogProduct := db.Model(&models.LogProduct{}).Where("date_txt=? AND an_product_id=?", dateNowTxt, v.AnProductId).FirstOrCreate(&setCreatedLogLogProduct)
		fmt.Println("role : ", createLogProduct.RowsAffected)

		resultData := map[string]interface{}{
			"no":            v.No,
			"an_product_id": v.AnProductId,
			"quantity":      v.Quantity,
			"product_code":  getProduct.ProductCode,
			"product_name":  getProduct.Name,
		}

		//update product stock
		if err := db.Table("product").Unscoped().Where("an_product_id=?", v.AnProductId).Debug().Updates(map[string]interface{}{
			"stock": remainingTodayStr,
		}).Error; err != nil {
			logrus.Info("update product error --> ", err.Error())
			return nil
		}

		notifyLine.SendMsgErrToLine(dateStart, dateEnd, getProduct.ProductCode, getProduct.Stock, quantityStr, orderManualAmountStr, remainingTodayStr)
		result = append(result, resultData)
	}

	logrus.Info("End Stock Cron Job", time.Now())
	return nil
}

func StockGetAll(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	var getProduct []models.Product = make([]models.Product, 0)
	product := db.Table("product").Unscoped()
	if err := product.Order("created_at desc").Find(&getProduct).Error; err != nil {
		logrus.Info("get product error --> ", err.Error())
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      getProduct,
	})
}

func CreateLogStock(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	//request
	payload := struct {
		ProductCode string `json:"product_code"`
		Name        string `json:"name"`
		Stock       string `json:"stock"`
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
	if err := product.Order("created_at desc").First(&getProduct).Error; err != nil {
		logrus.Info("get product error --> ", err.Error())
		return nil
	}

	stockPayloadConv, _ := strconv.Atoi(payload.Stock)
	stockConv, _ := strconv.Atoi(getProduct.Stock)

	sumStock := stockPayloadConv + stockConv
	sumStockStr := strconv.Itoa(sumStock)

	createLogStock := models.LogStock{
		ProductCode:    payload.ProductCode,
		ProductName:    getProduct.Name,
		StockIn:        payload.Stock,
		RemainingStock: getProduct.Stock,
		RemainingToday: sumStockStr,
	}
	db.Model(&models.LogStock{}).Create(&createLogStock)

	//update product stock
	if err := db.Table("product").Unscoped().Where("product_code=?", createLogStock.ProductCode).Debug().Updates(map[string]interface{}{
		"stock": sumStockStr,
	}).Error; err != nil {
		logrus.Info("update product error --> ", err.Error())
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
	})
}

func LogStockGetAll(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	var getLogStock []models.LogStock = make([]models.LogStock, 0)
	product := db.Model(&models.LogStock{})
	if err := product.Order("created_at desc").Find(&getLogStock).Error; err != nil {
		logrus.Info("get product error --> ", err.Error())
		return nil
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range getLogStock {

		resultData := map[string]interface{}{
			"created_at":      v.CreatedAt.Format("2006-01-02 15:04"),
			"product_code":    v.ProductCode,
			"product_name":    v.ProductName,
			"remaining_stock": v.RemainingStock,
			"stock_in":        v.StockIn,
			"remaining_today": v.RemainingToday,
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
