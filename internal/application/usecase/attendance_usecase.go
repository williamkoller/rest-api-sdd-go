package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

var ErrForbidden = errors.New("forbidden")

type AttendanceUseCase struct {
	attendanceRepo repository.AttendanceRepository
	userRepo       repository.UserRepository
	studentRepo    repository.StudentRepository
}

func NewAttendanceUseCase(attendanceRepo repository.AttendanceRepository, userRepo repository.UserRepository, studentRepo repository.StudentRepository) *AttendanceUseCase {
	return &AttendanceUseCase{attendanceRepo: attendanceRepo, userRepo: userRepo, studentRepo: studentRepo}
}

func (uc *AttendanceUseCase) BatchRecord(ctx context.Context, classID string, date time.Time, records []repository.AttendanceInput, teacherID, role string) (int, error) {
	if role == string(entity.RoleTeacher) {
		ok, err := uc.userRepo.IsTeacherOfClass(ctx, teacherID, classID)
		if err != nil {
			return 0, fmt.Errorf("attendance usecase: check teacher: %w", err)
		}
		if !ok {
			return 0, ErrForbidden
		}
	}

	count, err := uc.attendanceRepo.BatchUpsert(ctx, classID, date, records, teacherID)
	if err != nil {
		return 0, fmt.Errorf("attendance usecase: batch record: %w", err)
	}
	return count, nil
}

func (uc *AttendanceUseCase) GetByStudent(ctx context.Context, studentID string, from, to *time.Time, guardianID, role string) (*entity.AttendanceSummary, error) {
	if role == string(entity.RoleGuardian) {
		ok, err := uc.studentRepo.IsGuardianOf(ctx, guardianID, studentID)
		if err != nil {
			return nil, fmt.Errorf("attendance usecase: check guardian: %w", err)
		}
		if !ok {
			return nil, ErrForbidden
		}
	}

	summary, err := uc.attendanceRepo.FindByStudentID(ctx, studentID, from, to)
	if err != nil {
		return nil, fmt.Errorf("attendance usecase: get by student: %w", err)
	}
	return summary, nil
}
