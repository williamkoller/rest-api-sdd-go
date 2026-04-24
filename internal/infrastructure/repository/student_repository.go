package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	domainrepo "github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) domainrepo.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *entity.Student) error {
	m := model.StudentFromEntity(student)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("student repository: create: %w", err)
	}
	*student = *m.ToEntity()
	return nil
}

func (r *studentRepository) FindByID(ctx context.Context, id string) (*entity.Student, error) {
	var m model.StudentModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("student repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *studentRepository) FindByClassID(ctx context.Context, classID string) ([]*entity.Student, error) {
	var ms []model.StudentModel
	err := r.db.WithContext(ctx).
		Joins("JOIN enrollments ON enrollments.student_id = students.id").
		Where("enrollments.class_id = ? AND enrollments.status = 'active'", classID).
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("student repository: find by class id: %w", err)
	}
	students := make([]*entity.Student, len(ms))
	for i, m := range ms {
		m := m
		students[i] = m.ToEntity()
	}
	return students, nil
}

func (r *studentRepository) FindByCPF(ctx context.Context, schoolID, cpf string) (*entity.Student, error) {
	var m model.StudentModel
	if err := r.db.WithContext(ctx).Where("school_id = ? AND cpf = ?", schoolID, cpf).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("student repository: find by cpf: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *studentRepository) Update(ctx context.Context, student *entity.Student) error {
	m := model.StudentFromEntity(student)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("student repository: update: %w", err)
	}
	return nil
}

func (r *studentRepository) LinkGuardian(ctx context.Context, guardianID, studentID, relationship string) error {
	err := r.db.WithContext(ctx).Exec(
		"INSERT INTO guardian_students (guardian_id, student_id, relationship) VALUES (?, ?, ?) ON CONFLICT DO NOTHING",
		guardianID, studentID, relationship,
	).Error
	if err != nil {
		return fmt.Errorf("student repository: link guardian: %w", err)
	}
	return nil
}

func (r *studentRepository) IsGuardianOf(ctx context.Context, guardianID, studentID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("guardian_students").
		Where("guardian_id = ? AND student_id = ?", guardianID, studentID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("student repository: is guardian of: %w", err)
	}
	return count > 0, nil
}
