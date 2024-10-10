package converter

import (
	"github.com/NikolosHGW/auth/internal/infrastructure/db/repository/user/model"
	userpb "github.com/NikolosHGW/auth/pkg/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToGetResponseFromUserModel - функция, которая конвертирует модель в grpc response.
func ToGetResponseFromUserModel(user *model.User) *userpb.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &userpb.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      userpb.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
