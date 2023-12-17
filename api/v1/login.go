package v1

import (
	dbCon "api-medese/db"
	middleware "api-medese/middleware/jwt_user"
	"api-medese/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Login(ctx *fiber.Ctx) error {
	db := dbCon.DBConn
	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
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
			Message:   "invalid username and password",
			MessageTh: "โปรดระบุ username and password",
			Error:     "bad request",
		})
	}

	var user models.Users
	if err := db.Model(&models.Users{}).Where("username=?", payload.Username).First(&user).Error; err != nil {
		logrus.Errorln("get user : ", err)
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "invalid user",
			MessageTh: "ไม่มีผู้ใช้ในระบบ",
			Error:     "not found",
		})
	}

	token, err := middleware.CreateTokenUser(payload.Username)
	if err != nil {
		logrus.Errorln("CreateTokenUser : ", err)
		return ctx.Status(fiber.StatusNotFound).JSON(Result{
			Status:    fiber.StatusNotFound,
			Message:   "create token error",
			MessageTh: "create token error",
			Error:     "not found",
		})
	}

	result := map[string]interface{}{
		"token": token,
	}

	return ctx.Status(fiber.StatusOK).JSON(Result{
		Status:    fiber.StatusOK,
		Message:   "success",
		MessageTh: "สำเร็จ",
		Data:      result,
	})

}
