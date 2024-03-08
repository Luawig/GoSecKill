package main

import (
	"GoSecKill/internal/config"
	"GoSecKill/internal/database"
	"GoSecKill/internal/routers"
	"GoSecKill/pkg/log"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Load application configuration
	if err := config.LoadConfig("./config"); err != nil {
		panic(err)
	}

	// Initialize logger
	log.InitLogger()
	zap.L().Info("log init success")

	// Initialize database
	db := database.InitDB()

	// Initialize the web application
	app := iris.New()
	app.Logger().SetLevel("debug")

	// Register the view engine
	template := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	// Register the routes
	app.HandleDir("/assets", "./web/assets")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "There was something wrong with the request!"))
		ctx.ViewData("status", ctx.GetStatusCode())
		ctx.ViewLayout("")
		_ = ctx.View("shared/error.html")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register the routes
	routers.InitRoutes(app, db, ctx)

	// Register the product routes
	//productRepository := repository.NewProductRepository(db)
	//productService := services.NewProductService(productRepository)
	//productParty := app.Party("/product")
	//product := mvc.New(productParty)
	//product.Register(ctx, productService)
	//product.Handle(new(controllers.ProductController))

	// Start the web application
	err := app.Run(
		iris.Addr(viper.GetString("server.port")),
		iris.WithCharset("UTF-8"),
		iris.WithOptimizations,
		iris.WithoutServerError(iris.ErrServerClosed),
	)
	if err != nil {
		zap.L().Fatal("app run failed", zap.Error(err))
	}
}
