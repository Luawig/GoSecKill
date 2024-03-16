package main

import (
	"GoSecKill/internal/config"
	"GoSecKill/internal/database"
	"GoSecKill/internal/routers"
	"GoSecKill/pkg/log"
	"GoSecKill/pkg/mq"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
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
	template := iris.HTML("./web/server/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	// Register the routes
	app.HandleDir("/html", "./web/server/htmlProductShow")
	app.HandleDir("/assets", "./web/server/assets")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "There was something wrong with the request!"))
		ctx.ViewData("status", ctx.GetStatusCode())
		ctx.ViewLayout("")
		_ = ctx.View("shared/error.html")
	})

	// Initialize the message queue
	rabbitmq := mq.NewRabbitMQSimple("go_seckill")

	session := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * 60 * 60,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register the routes
	routers.InitServerRoutes(app, db, ctx, session, rabbitmq)

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
