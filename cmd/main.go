package main

import (
	"context"
	"log"

	"github.com/NikolosHGW/auth/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("ошибка инициализации приложения: %s", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("ошибка при запуске приложения: %s", err.Error())
	}
}
