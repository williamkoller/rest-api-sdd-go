package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *entity.Enrollment) error
	FindByID(ctx context.Context, id string) (*entity.Enrollment, error)
	FindActiveByStudentAndYear(ctx context.Context, studentID string, year int) (*entity.Enrollment, error)
	FindByClassID(ctx context.Context, classID string) ([]*entity.Enrollment, error)
	FindByStudentID(ctx context.Context, studentID string) ([]*entity.Enrollment, error)
	Unenroll(ctx context.Context, enrollmentID string) error
}
