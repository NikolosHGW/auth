package converter

import (
	serviceUser "github.com/NikolosHGW/auth/internal/service/user/model"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// APICreateUserToServiceUser конвертирует proto CreateRequest в модель сервиса юзера.
func APICreateUserToServiceUser(apiUser *userpb.CreateRequest) *serviceUser.User {
	return &serviceUser.User{
		Name:     apiUser.Name,
		Email:    apiUser.Email,
		Password: apiUser.Password,
		Role:     int32(apiUser.Role),
	}
}

// ServiceUserToAPIGetUser конвертирует модель сервиса юзера в proto GetResponse.
func ServiceUserToAPIGetUser(serviceUser *serviceUser.User) *userpb.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if serviceUser.UpdatedAt.Valid {
		updatedAt = timestamppb.New(serviceUser.UpdatedAt.Time)
	}

	return &userpb.GetResponse{
		Id:        serviceUser.ID,
		Name:      serviceUser.Name,
		Email:     serviceUser.Email,
		Role:      userpb.Role(serviceUser.Role),
		CreatedAt: timestamppb.New(serviceUser.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// APIUpdateUserToServiceUser конвертирует proto UpdateRequest в модель сервиса юзера.
func APIUpdateUserToServiceUser(apiUser *userpb.UpdateRequest) *serviceUser.User {
	return &serviceUser.User{
		ID:    apiUser.Id,
		Name:  apiUser.Name.GetValue(),
		Email: apiUser.Email.GetValue(),
	}
}
