package repository

import (
	"context"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type GradeInput struct {
	StudentID string
	Value     float64
}

type GradeFilters struct {
	Subject      string
	Period       string
	AcademicYear *int
}

type GradeRepository interface {
	BatchUpsert(ctx context.Context, classID, subject, period string, grades []GradeInput, recordedBy string) (int, error)
	FindByStudentID(ctx context.Context, studentID string, filters GradeFilters) ([]*entity.Grade, error)
}
