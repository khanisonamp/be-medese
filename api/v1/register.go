package v1

import (
	dbCon "api-medese/db"
	"api-medese/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	helper "api-medese/helper/password"
)

func RegisterMember(ctx *fiber.Ctx) error {
	db := dbCon.DBConn

	payload := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		logrus.Errorln("BodyParser : ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid body",
			MessageTh: "body ไม่ถูกต้อง",
			Error:     "bad request",
		})
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		logrus.Errorln("validate : ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(Result{
			Status:    fiber.StatusBadRequest,
			Message:   "invalid request body",
			MessageTh: "โปรดระบุ request body",
			Error:     "bad request",
		})
	}

	pass := helper.GenPassword()
	hashPassword, _ := helper.HashPassword(pass)

	setCreatedUser := models.Users{
		Username:     payload.Username,
		Password:     pass,
		PasswordHash: hashPassword,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
	}

	createUSer := db.Model(&models.Users{}).Where("username=?", payload.Username).FirstOrCreate(&setCreatedUser)
	if createUSer.RowsAffected == 0 {
		return ctx.Status(fiber.StatusConflict).JSON(Result{
			Status:    fiber.StatusConflict,
			Message:   "Duplicate User",
			MessageTh: "Duplicate User",
			Error:     "Conflict",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      setCreatedUser,
	})
}
