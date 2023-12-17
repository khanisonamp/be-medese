package fiber

import (
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type FiberContext struct {
	*fiber.Ctx
}

type FiberContextJwt struct {
	*fiber.Ctx
	err error
}

func NewFiberContext(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

func NewFiberContextJwt(c *fiber.Ctx, err error) *FiberContextJwt {
	return &FiberContextJwt{Ctx: c, err: err}
}

func NewRouter(c *fiber.Ctx, err error) *FiberContextJwt {
	return &FiberContextJwt{Ctx: c, err: err}
}

// composition over inheritance
type Context interface {
	BodyParser(interface{}) error
	SendResponse(int, interface{})
	ClaimsJWT() *jwt.Token
	Query(string) string
	Get(string) string
	SendString(string) error
	Params(key string) string
	FormFile(key string) (*multipart.FileHeader, error)
	Map(v map[string]interface{}) map[string]interface{}
	Locals(key string) interface{}
	Next() error
	RequestBody() []byte
}

func (c *FiberContext) RequestBody() []byte {
	return c.Ctx.Request().Body()
}

func (c *FiberContext) SendString(body string) error {
	return c.Ctx.SendString(body)
}

func (c *FiberContext) BodyParser(body interface{}) error {
	return c.Ctx.BodyParser(body)
}

func (c *FiberContext) SendResponse(code int, res interface{}) {
	c.Status(code).JSON(res)
	return
}

func (c *FiberContext) ClaimsJWT() *jwt.Token {
	return c.Ctx.Locals("user").(*jwt.Token)
}

func (c *FiberContext) Query(key string) string {
	return strings.TrimSpace(c.Ctx.Query(key))
}

func (c *FiberContext) Get(key string) string {
	return c.Ctx.Get(key)
}

func (c *FiberContext) Params(key string) string {
	return strings.TrimSpace(c.Ctx.Params(key))
}

func (c *FiberContext) FormFile(key string) (*multipart.FileHeader, error) {
	return c.Ctx.FormFile(key)
}

func (c *FiberContext) Map(v map[string]interface{}) map[string]interface{} {
	return v
}

func (c *FiberContext) Locals(key string) interface{} {
	return c.Ctx.Locals(key)
}

func (c *FiberContext) Next() error {
	return c.Ctx.Next()
}
