package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"github.com/NikolosHGW/platform-common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App - структура приложения.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("ошибка запуска gPRC сервера: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("ошибка запуска HTTP-сервера: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := userpb.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().GetRunAddress(), opts)
	if err != nil {
		return fmt.Errorf("ошибка при инициализации HTTP-сервера: %w", err)
	}

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().HTTPRunAddress(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,   // Таймаут на чтение заголовков
		ReadTimeout:       10 * time.Second,  // (опционально) Таймаут на чтение всего запроса
		WriteTimeout:      10 * time.Second,  // (опционально) Таймаут на запись ответа
		IdleTimeout:       120 * time.Second, // (опционально) Таймаут для неактивных соединений
	}

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
				a.serviceProvider.GRPCConfig().GetRunAddress(),
			)
			a.grpcServer.GracefulStop()

			return nil
		},
	)

	log.Printf("Запуск GRPC сервера на адресе %s ...", a.serviceProvider.GRPCConfig().GetRunAddress())

	if err := a.grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %w", err)
	}

	return nil
}

func (a *App) runHTTPServer() error {
	closer.Add(
		func() error {
			log.Printf(
				"Получен сигнал завершения, отключаем сервер по адресу %s ...",
				a.serviceProvider.HTTPConfig().HTTPRunAddress(),
			)

			ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancelShutDown()

			return a.httpServer.Shutdown(ctxShutDown)
		},
	)
	log.Printf("Запуск HTTP-сервера на адресе %s ...", a.serviceProvider.HTTPConfig().HTTPRunAddress())

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ошибка при запуске HTTP-сервера: %w", err)
	}

	return nil
}
