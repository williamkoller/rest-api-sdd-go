package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type StudentRepository interface {
	Create(ctx context.Context, student *entity.Student) error
	FindByID(ctx context.Context, id string) (*entity.Student, error)
	FindByClassID(ctx context.Context, classID string) ([]*entity.Student, error)
	FindByCPF(ctx context.Context, schoolID, cpf string) (*entity.Student, error)
	Update(ctx context.Context, student *entity.Student) error
	LinkGuardian(ctx context.Context, guardianID, studentID, relationship string) error
	IsGuardianOf(ctx context.Context, guardianID, studentID string) (bool, error)
}
