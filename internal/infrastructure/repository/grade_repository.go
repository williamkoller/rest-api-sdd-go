package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) domainrepo.GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) BatchUpsert(ctx context.Context, classID, subject, period string, grades []domainrepo.GradeInput, recordedBy string) (int, error) {
	type enrollmentRow struct {
		ID        string
		StudentID string
	}
	var enrollments []enrollmentRow
	err := r.db.WithContext(ctx).Table("enrollments").
		Select("id, student_id").
		Where("class_id = ? AND status = 'active'", classID).
		Scan(&enrollments).Error
	if err != nil {
		return 0, fmt.Errorf("grade repository: fetch enrollments: %w", err)
	}

	studentToEnrollment := make(map[string]string, len(enrollments))
	for _, e := range enrollments {
		studentToEnrollment[e.StudentID] = e.ID
	}

	var ms []model.GradeModel
	for _, g := range grades {
		enrollmentID, ok := studentToEnrollment[g.StudentID]
		if !ok {
			continue
		}
		ms = append(ms, model.GradeModel{
			ID:           uuid.New().String(),
			EnrollmentID: enrollmentID,
			Subject:      subject,
			Period:       period,
			Value:        g.Value,
			RecordedBy:   recordedBy,
		})
	}

	if len(ms) == 0 {
		return 0, nil
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "enrollment_id"}, {Name: "subject"}, {Name: "period"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "recorded_by", "updated_at"}),
		}).
		Create(&ms).Error; err != nil {
		return 0, fmt.Errorf("grade repository: batch upsert: %w", err)
	}
	return len(ms), nil
}

func (r *gradeRepository) FindByStudentID(ctx context.Context, studentID string, filters domainrepo.GradeFilters) ([]*entity.Grade, error) {
	var ms []model.GradeModel
	q := r.db.WithContext(ctx).
		Joins("JOIN enrollments ON enrollments.id = grades.enrollment_id").
		Where("enrollments.student_id = ?", studentID)
	if filters.Subject != "" {
		q = q.Where("grades.subject = ?", filters.Subject)
	}
	if filters.Period != "" {
		q = q.Where("grades.period = ?", filters.Period)
	}
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("grade repository: find by student id: %w", err)
	}
	result := make([]*entity.Grade, len(ms))
	for i, m := range ms {
		m := m
		result[i] = m.ToEntity()
	}
	return result, nil
}
