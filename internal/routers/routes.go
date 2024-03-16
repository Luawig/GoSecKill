package routers

import (
	"GoSecKill/internal/middleware"
	"GoSecKill/internal/services"
	"GoSecKill/pkg/mq"
	"GoSecKill/pkg/repositories"
	controllers2 "GoSecKill/web/admin/controllers"
	"GoSecKill/web/server/controllers"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gorm.io/gorm"
)

func InitAdminRoutes(app *iris.Application, db *gorm.DB, ctx context.Context) {
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers2.ProductController))
	product.Handle(controllers2.NewProductController(productService))

	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers2.OrderController))
	order.Handle(controllers2.NewOrderController(orderService))
}

func InitServerRoutes(app *iris.Application, db *gorm.DB, ctx context.Context, sessions *sessions.Sessions, rabbitmq *mq.RabbitMQ) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))
	user.Handle(controllers.NewUserController(userService, sessions))

	product := repositories.NewProductRepository(db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService, rabbitmq)
	pro.Handle(controllers.NewProductController(productService, orderService, sessions))
}
