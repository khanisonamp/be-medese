package router

import (
	api "api-medese/api/v1"
	contextFiber "api-medese/framework/fiber"
)

func SetRouter(app *contextFiber.FiberApp) {
	backend := app.Group("/ware-house")
	v1 := backend.Group("/api/v1")

	user := v1.Group("/user")
	user.Post("/register", api.RegisterMember)
	user.Get("/", api.GetUserFirst)

	login := v1.Group("/login")
	login.Group("/", api.Login)

	order := v1.Group("/order")
	order.Post("/", api.OrderGetAll)

	orderProduct := v1.Group("/order-product")
	orderProduct.Post("/dashboard", api.OrderProductDashBoard)

	logProduct := v1.Group("/log-product")
	logProduct.Get("/", api.ScriptCrestedLogStockAuto)
	logProduct.Get("/month", api.GetLogProduct)

	product := v1.Group("/product")
	product.Get("/", api.StockGetAll)
	product.Post("/", api.CreateLogStock)
	product.Get("/log-create-stock", api.LogStockGetAll)

}
