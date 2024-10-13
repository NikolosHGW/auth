package app

import (
	"context"
	"log"

	apiUser "github.com/NikolosHGW/auth/internal/api/user"
	"github.com/NikolosHGW/auth/internal/closer"
	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository"
	repositoryUser "github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user"
	"github.com/NikolosHGW/auth/internal/service"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	pgConfig   config.PGConfiger
	grpcConfig config.GRPCConfiger

	pgxPool *pgxpool.Pool

	userRepo repository.UserRepository

	userService service.UserService

	userImplementation *apiUser.Implementation
}

func newServiceProvider() *serviceProvider {
	return new(serviceProvider)
}

// PGConfig - синглтон для постгрес конфига.
func (s *serviceProvider) PGConfig() config.PGConfiger {
	if s.pgConfig == nil {
		pgConfig, err := config.NewPG()
		if err != nil {
			log.Fatalf("ошибка при инициализации объекта pgConfig: %s", err.Error())
		}

		s.pgConfig = pgConfig
	}

	return s.pgConfig
}

// GRPCConfig - синглтон для grpc конфига.
func (s *serviceProvider) GRPCConfig() config.GRPCConfiger {
	if s.grpcConfig == nil {
		grpcConfig, err := config.NewGRPC()
		if err != nil {
			log.Fatalf("ошибка при инициализации объекта grpcConfig: %s", err.Error())
		}

		s.grpcConfig = grpcConfig
	}

	return s.grpcConfig
}

// PGXPool - синглтон для pgx пула.
func (s *serviceProvider) PGXPool(ctx context.Context) *pgxpool.Pool {
	if s.pgxPool == nil {
		pgxPool, err := pgxpool.New(ctx, s.PGConfig().GetDatabaseDSN())
		if err != nil {
			log.Fatalf("ошибка при инициализации pgx pool: %s", err.Error())
		}

		err = pgxPool.Ping(ctx)
		if err != nil {
			log.Fatalf("ошибка во время пинга к бд: %s", err.Error())
		}

		closer.Add(
			func() error {
				pgxPool.Close()

				return nil
			},
		)

		s.pgxPool = pgxPool
	}

	return s.pgxPool
}

// UserRepository - синглтон для репозитория юзера.
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = repositoryUser.NewUser(s.PGXPool(ctx))
	}

	return s.userRepo
}

// UserService - синглтон для сервиса юзера.
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = serviceUser.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

// UserImplementation - синглтон для grpc сервера юзера.
func (s *serviceProvider) UserImplementation(ctx context.Context) *apiUser.Implementation {
	if s.userImplementation == nil {
		s.userImplementation = apiUser.NewImplementation(s.UserService(ctx))
	}

	return s.userImplementation
}
