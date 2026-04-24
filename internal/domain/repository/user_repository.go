package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	FindTeachersByClassID(ctx context.Context, classID string) ([]*entity.User, error)
	IsTeacherOfClass(ctx context.Context, teacherID, classID string) (bool, error)
}
