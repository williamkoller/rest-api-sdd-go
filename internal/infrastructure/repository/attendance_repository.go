package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) domainrepo.AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) BatchUpsert(ctx context.Context, classID string, date time.Time, records []domainrepo.AttendanceInput, recordedBy string) (int, error) {
	// Get enrollment IDs for each student in the class
	type enrollmentRow struct {
		ID        string
		StudentID string
	}
	var enrollments []enrollmentRow
	err := r.db.WithContext(ctx).
		Table("enrollments").
		Select("id, student_id").
		Where("class_id = ? AND status = 'active'", classID).
		Scan(&enrollments).Error
	if err != nil {
		return 0, fmt.Errorf("attendance repository: fetch enrollments: %w", err)
	}

	studentToEnrollment := make(map[string]string, len(enrollments))
	for _, e := range enrollments {
		studentToEnrollment[e.StudentID] = e.ID
	}

	var ms []model.AttendanceRecordModel
	for _, rec := range records {
		enrollmentID, ok := studentToEnrollment[rec.StudentID]
		if !ok {
			continue
		}
		ms = append(ms, model.AttendanceRecordModel{
			ID:           uuid.New().String(),
			EnrollmentID: enrollmentID,
			Date:         date,
			Status:       string(rec.Status),
			Note:         rec.Note,
			RecordedBy:   recordedBy,
		})
	}

	if len(ms) == 0 {
		return 0, nil
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "enrollment_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "note", "recorded_by"}),
		}).
		Create(&ms).Error; err != nil {
		return 0, fmt.Errorf("attendance repository: batch upsert: %w", err)
	}
	return len(ms), nil
}

func (r *attendanceRepository) FindByStudentID(ctx context.Context, studentID string, from, to *time.Time) (*entity.AttendanceSummary, error) {
	var ms []model.AttendanceRecordModel
	q := r.db.WithContext(ctx).
		Joins("JOIN enrollments ON enrollments.id = attendance_records.enrollment_id").
		Where("enrollments.student_id = ?", studentID)
	if from != nil {
		q = q.Where("attendance_records.date >= ?", from)
	}
	if to != nil {
		q = q.Where("attendance_records.date <= ?", to)
	}
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("attendance repository: find by student id: %w", err)
	}

	records := make([]*entity.AttendanceRecord, len(ms))
	present := 0
	for i, m := range ms {
		m := m
		records[i] = m.ToEntity()
		if m.Status == "present" {
			present++
		}
	}
	total := len(records)
	rate := 0.0
	if total > 0 {
		rate = float64(present) / float64(total)
	}
	return &entity.AttendanceSummary{
		Records: records,
		Rate:    rate,
		Total:   total,
		Present: present,
	}, nil
}
