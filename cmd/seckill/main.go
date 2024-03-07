package main

import (
	"GoSecKill/internal/config"
	"GoSecKill/pkg/log"

	"go.uber.org/zap"
)

func main() {
	if err := config.LoadConfig("./config"); err != nil {
		panic(err)
	}

	log.InitLogger()
	zap.L().Info("log init success")
}
