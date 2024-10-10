package converter

import (
	repoModel "github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user/model"
	serviceModel "github.com/NikolosHGW/auth/internal/service/user/model"
)

// UserServiceToUserRepo конвертирует модель сервиса юзера в модель репы юзера.
func UserServiceToUserRepo(serviceUser *serviceModel.User) *repoModel.User {
	return &repoModel.User{
		ID:        serviceUser.ID,
		Name:      serviceUser.Name,
		Email:     serviceUser.Email,
		Password:  serviceUser.Password,
		Role:      serviceUser.Role,
		CreatedAt: serviceUser.CreatedAt,
		UpdatedAt: serviceUser.UpdatedAt,
	}
}

// UserRepoToServiceRepo конвертирует модель репы юзера в модель сервиса юзера.
func UserRepoToServiceRepo(repoUser *repoModel.User) *serviceModel.User {
	return &serviceModel.User{
		ID:        repoUser.ID,
		Name:      repoUser.Name,
		Email:     repoUser.Email,
		Password:  repoUser.Password,
		Role:      repoUser.Role,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}
}
