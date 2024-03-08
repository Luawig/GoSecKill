package routers

import (
	"GoSecKill/internal/controllers"
	"GoSecKill/internal/services"
	"GoSecKill/pkg/repositories"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
)

func InitRoutes(app *iris.Application, db *gorm.DB, ctx context.Context) {
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))
	product.Handle(controllers.NewProductController(productService))

	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	order.Handle(controllers.NewOrderController(orderService))
}
