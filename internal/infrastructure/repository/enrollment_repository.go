package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) domainrepo.EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	m := model.EnrollmentFromEntity(enrollment)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("enrollment repository: create: %w", err)
	}
	*enrollment = *m.ToEntity()
	return nil
}

func (r *enrollmentRepository) FindByID(ctx context.Context, id string) (*entity.Enrollment, error) {
	var m model.EnrollmentModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("enrollment repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *enrollmentRepository) FindActiveByStudentAndYear(ctx context.Context, studentID string, year int) (*entity.Enrollment, error) {
	var m model.EnrollmentModel
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND academic_year = ? AND status = 'active'", studentID, year).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("enrollment repository: find active by student and year: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *enrollmentRepository) FindByClassID(ctx context.Context, classID string) ([]*entity.Enrollment, error) {
	var ms []model.EnrollmentModel
	if err := r.db.WithContext(ctx).Where("class_id = ? AND status = 'active'", classID).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("enrollment repository: find by class id: %w", err)
	}
	enrollments := make([]*entity.Enrollment, len(ms))
	for i, m := range ms {
		m := m
		enrollments[i] = m.ToEntity()
	}
	return enrollments, nil
}

func (r *enrollmentRepository) FindByStudentID(ctx context.Context, studentID string) ([]*entity.Enrollment, error) {
	var ms []model.EnrollmentModel
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("enrollment repository: find by student id: %w", err)
	}
	enrollments := make([]*entity.Enrollment, len(ms))
	for i, m := range ms {
		m := m
		enrollments[i] = m.ToEntity()
	}
	return enrollments, nil
}

func (r *enrollmentRepository) Unenroll(ctx context.Context, enrollmentID string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&model.EnrollmentModel{}).
		Where("id = ?", enrollmentID).
		Updates(map[string]any{
			"status":        "unenrolled",
			"unenrolled_at": now,
		}).Error
	if err != nil {
		return fmt.Errorf("enrollment repository: unenroll: %w", err)
	}
	return nil
}
