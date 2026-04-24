package repository

import (
	"context"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type AttendanceInput struct {
	StudentID string
	Status    entity.AttendanceStatus
	Note      string
}

type AttendanceRepository interface {
	BatchUpsert(ctx context.Context, classID string, date time.Time, records []AttendanceInput, recordedBy string) (int, error)
	FindByStudentID(ctx context.Context, studentID string, from, to *time.Time) (*entity.AttendanceSummary, error)
}
