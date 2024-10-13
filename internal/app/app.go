package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/NikolosHGW/auth/internal/closer"
	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp - конструктор Приложения App, где происходит вся инициализация приложения.
func NewApp(ctx context.Context) (*App, error) {
	a := new(App)

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при инициализации всех компонентов: %w", err)
	}

	return a, nil
}

// Run запускает сервер
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
		log.Println("Приложение завершено корректно.")
	}()

	if err := a.runGRPCServer(); err != nil {
		return fmt.Errorf("ошибка запуска gPRC сервера: %w", err)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("не удалось инициализировать компонент: %w", err)
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка при подгрузки конфигурации приложения: %w", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	userpb.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	listen, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().GetRunAddress())
	if err != nil {
		return fmt.Errorf("не удалось прослушать TCP: %w", err)
	}

	closer.Add(
		func() error {
			log.Printf(
				"Получен сигнал завершения, отключаем сервер по адресу %s ...",
				a.serviceProvider.grpcConfig.GetRunAddress(),
			)
			a.grpcServer.GracefulStop()

			return nil
		},
	)

	log.Printf("Запуск GRPC сервера на адресе %s", a.serviceProvider.GRPCConfig().GetRunAddress())

	if err := a.grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %w", err)
	}

	return nil
}
