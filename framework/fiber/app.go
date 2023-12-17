package fiber

import (
	"fmt"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberApp struct {
	*fiber.App
}

func NewFiberApp() *FiberApp {
	app := fiber.New(fiber.Config{
		BodyLimit: 25 * 1024 * 1024,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := ctx.Context().Response.StatusCode()

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			fmt.Println(err)
			// Send custom error page
			if err != nil {

				if code == 404 {
					return ctx.Status(code).JSON(fiber.Map{
						"status":     "404",
						"message_th": "ไม่พบข้อมูล",
						"message_en": "Not Found",
					})
				}

				if code == 200 {
					// In case the SendFile fails
					return ctx.Status(500).JSON(fiber.Map{
						"status":     "500",
						"message_en": "Internal server error",
						"message_th": "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์",
					})
				}

			}

			// Return from handler
			return nil

		},
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(recover.New(recover.ConfigDefault))

	app.Use(logger.New(logger.Config{
		Format:     "${blue}${time} ${yellow}${status} - ${red}${latency} ${cyan}${method} ${path} ${green} ${ip} ${ua} ${reset}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
		Output:     os.Stdout,
	}))

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Status(200).SendString("Api i-claim-v2")

	})

	app.Get("/backend", func(c *fiber.Ctx) error {

		return c.Status(200).SendString("Hello World...")

	})

	app.Get("/backend/docs/*", swagger.New(swagger.Config{
		DeepLinking: true,
		URL:         "doc.json",
	}))
	return &FiberApp{app}
}

func (f *FiberApp) POST(path string, handler func(Context)) {
	f.App.Post(path, func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	})

}

func (f *FiberApp) GET(path string, handler func(Context)) {
	f.App.Get(path, func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	})

}

func (f *FiberApp) DELETE(path string, handler func(Context)) {
	f.App.Delete(path, func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	})

}

func (f *FiberApp) PUT(path string, handler func(Context)) {
	f.App.Put(path, func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	})

}

func (f *FiberApp) GROUP(path string, handler func(Context)) {
	f.App.Put(path, func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	})
}

func NewFiberHandler(handler func(Context)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler(NewFiberContext(c))
		return nil
	}
}
