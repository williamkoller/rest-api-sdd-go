package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type StudentUseCase struct {
	studentRepo    repository.StudentRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewStudentUseCase(studentRepo repository.StudentRepository, enrollmentRepo repository.EnrollmentRepository) *StudentUseCase {
	return &StudentUseCase{studentRepo: studentRepo, enrollmentRepo: enrollmentRepo}
}

func (uc *StudentUseCase) Enroll(ctx context.Context, classID, schoolID string, req EnrollRequest) (*entity.Enrollment, error) {
	var student *entity.Student

	if req.StudentID != "" {
		var err error
		student, err = uc.studentRepo.FindByID(ctx, req.StudentID)
		if err != nil {
			return nil, fmt.Errorf("student usecase: enroll find student: %w", err)
		}
		if student == nil {
			return nil, ErrNotFound
		}
	} else {
		// Create new student
		regNum := uuid.New().String()[:8]
		student = &entity.Student{
			ID:                 uuid.New().String(),
			SchoolID:           schoolID,
			Name:               req.Name,
			BirthDate:          req.BirthDate,
			CPF:                req.CPF,
			RegistrationNumber: regNum,
			Active:             true,
		}
		if err := uc.studentRepo.Create(ctx, student); err != nil {
			return nil, fmt.Errorf("student usecase: create student: %w", err)
		}
	}

	enrollment := &entity.Enrollment{
		ID:           uuid.New().String(),
		StudentID:    student.ID,
		ClassID:      classID,
		AcademicYear: req.AcademicYear,
		EnrolledAt:   req.EnrolledAt,
		Status:       entity.EnrollmentStatusActive,
	}
	if err := uc.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, fmt.Errorf("student usecase: create enrollment: %w", err)
	}
	return enrollment, nil
}

func (uc *StudentUseCase) FindByClassID(ctx context.Context, classID string) ([]*entity.Student, error) {
	students, err := uc.studentRepo.FindByClassID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("student usecase: find by class id: %w", err)
	}
	return students, nil
}

func (uc *StudentUseCase) FindByID(ctx context.Context, id string) (*entity.Student, error) {
	student, err := uc.studentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("student usecase: find by id: %w", err)
	}
	if student == nil {
		return nil, ErrNotFound
	}
	return student, nil
}

func (uc *StudentUseCase) Update(ctx context.Context, id string, updates map[string]any) (*entity.Student, error) {
	student, err := uc.studentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("student usecase: update: %w", err)
	}
	if student == nil {
		return nil, ErrNotFound
	}
	if v, ok := updates["name"].(string); ok {
		student.Name = v
	}
	if v, ok := updates["cpf"].(string); ok {
		student.CPF = v
	}
	if err := uc.studentRepo.Update(ctx, student); err != nil {
		return nil, fmt.Errorf("student usecase: update: %w", err)
	}
	return student, nil
}

func (uc *StudentUseCase) Unenroll(ctx context.Context, classID, studentID string) error {
	enrollments, err := uc.enrollmentRepo.FindByClassID(ctx, classID)
	if err != nil {
		return fmt.Errorf("student usecase: unenroll: %w", err)
	}
	for _, e := range enrollments {
		if e.StudentID == studentID {
			return uc.enrollmentRepo.Unenroll(ctx, e.ID)
		}
	}
	return ErrNotFound
}

type EnrollRequest struct {
	StudentID    string
	Name         string
	BirthDate    time.Time
	CPF          string
	AcademicYear int
	EnrolledAt   time.Time
}
