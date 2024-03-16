package main

import (
	"GoSecKill/internal/database"
	"GoSecKill/internal/services"
	"GoSecKill/pkg/mq"
	"GoSecKill/pkg/repositories"
)

func main() {

	db := database.GetDB()

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)

	rabbitmqConsumer := mq.NewRabbitMQSimple("go_seckill")
	rabbitmqConsumer.ConsumeSimple(orderService, productService)
}
