package usecase

import (
	"context"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
)

type GradeUseCase struct {
	gradeRepo   repository.GradeRepository
	userRepo    repository.UserRepository
	studentRepo repository.StudentRepository
}

func NewGradeUseCase(gradeRepo repository.GradeRepository, userRepo repository.UserRepository, studentRepo repository.StudentRepository) *GradeUseCase {
	return &GradeUseCase{gradeRepo: gradeRepo, userRepo: userRepo, studentRepo: studentRepo}
}

func (uc *GradeUseCase) BatchUpsert(ctx context.Context, classID, subject, period string, grades []repository.GradeInput, teacherID, role string) (int, error) {
	if role == string(entity.RoleTeacher) {
		ok, err := uc.userRepo.IsTeacherOfClass(ctx, teacherID, classID)
		if err != nil {
			return 0, fmt.Errorf("grade usecase: check teacher: %w", err)
		}
		if !ok {
			return 0, ErrForbidden
		}
	}

	count, err := uc.gradeRepo.BatchUpsert(ctx, classID, subject, period, grades, teacherID)
	if err != nil {
		return 0, fmt.Errorf("grade usecase: batch upsert: %w", err)
	}
	return count, nil
}

func (uc *GradeUseCase) GetByStudent(ctx context.Context, studentID string, filters repository.GradeFilters, guardianID, role string) ([]*entity.Grade, error) {
	if role == string(entity.RoleGuardian) {
		ok, err := uc.studentRepo.IsGuardianOf(ctx, guardianID, studentID)
		if err != nil {
			return nil, fmt.Errorf("grade usecase: check guardian: %w", err)
		}
		if !ok {
			return nil, ErrForbidden
		}
	}

	grades, err := uc.gradeRepo.FindByStudentID(ctx, studentID, filters)
	if err != nil {
		return nil, fmt.Errorf("grade usecase: get by student: %w", err)
	}
	return grades, nil
}
