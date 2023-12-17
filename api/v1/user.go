package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	helper "api-medese/helper"
)

func Register(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid body",
			MessageTh: "body ไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "validate body",
			MessageTh: "body ไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	setModelCreate := models.Users{
		Username: payload.Username,
	}

	if err := db.Model(&models.Users{}).Create(&setModelCreate).Debug().Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "create user error",
			MessageTh: "create user error",
			Error:     "bad request",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      payload,
	})

}

func GetUserFirst(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	rstart := "2023-11-20 14:44"
	rend := "2023-11-21 14:45"

	start, _ := time.Parse("2006-01-02 15:04", rstart)
	end, _ := time.Parse("2006-01-02 15:04", rend)

	startD := helper.ParseLocal(start)
	endD := helper.ParseLocal(end)

	fmt.Println("startD parse ===> ", startD)
	fmt.Println("endD parse ===> ", endD)

	var getOrderProduct []models.OrderProduct = make([]models.OrderProduct, 0)
	if err := db.Table("order_product").Unscoped().Where("created_at BETWEEN ? AND ?", startD, endD).Debug().Find(&getOrderProduct).Order("created_at DESC").Error; err != nil {
		logrus.Info("error --> ", err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "get order_product error",
			MessageTh: "get order_product error",
			Error:     "not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      getOrderProduct,
	})

}
