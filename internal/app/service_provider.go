package app

import (
	"context"
	"fmt"
	"log"

	apiUser "github.com/NikolosHGW/auth/internal/api/user"
	"github.com/NikolosHGW/auth/internal/infrastructure/config"
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository"
	repositoryUser "github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user"
	"github.com/NikolosHGW/auth/internal/service"
	serviceUser "github.com/NikolosHGW/auth/internal/service/user"
	"github.com/NikolosHGW/platform-common/pkg/closer"
	"github.com/NikolosHGW/platform-common/pkg/db"
	"github.com/NikolosHGW/platform-common/pkg/db/pg"
	"github.com/NikolosHGW/platform-common/pkg/db/transaction"
)

type serviceProvider struct {
	pgConfig   config.PGConfiger
	grpcConfig config.GRPCConfiger

	dbClient  db.Client
	txManager db.TxManager

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

// DBClient - синглтон для бд.
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		dbClient, err := pg.New(ctx, s.PGConfig().GetDatabaseDSN())
		if err != nil {
			log.Fatalf("ошибка при инициализации клиента бд: %s", err.Error())
		}

		err = dbClient.DB().PingContext(ctx)
		if err != nil {
			log.Fatalf("ошибка во время пинга к бд: %s", err.Error())
		}

		closer.Add(
			func() error {
				err := dbClient.Close()
				if err != nil {
					return fmt.Errorf("ошибка при закрытии БД клиента: %w", err)
				}

				return nil
			},
		)

		s.dbClient = dbClient
	}

	return s.dbClient
}

func (s *serviceProvider) TXManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository - синглтон для репозитория юзера.
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = repositoryUser.NewUser(s.DBClient(ctx))
	}

	return s.userRepo
}

// UserService - синглтон для сервиса юзера.
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = serviceUser.NewService(s.UserRepository(ctx), s.TXManager(ctx))
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
