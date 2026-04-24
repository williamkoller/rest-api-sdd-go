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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainrepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := model.UserFromEntity(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("user repository: create: %w", err)
	}
	*user = *m.ToEntity()
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("user repository: find by id: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("user repository: find by email: %w", err)
	}
	return m.ToEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	m := model.UserFromEntity(user)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("user repository: update: %w", err)
	}
	return nil
}

func (r *userRepository) FindTeachersByClassID(ctx context.Context, classID string) ([]*entity.User, error) {
	var ms []model.UserModel
	err := r.db.WithContext(ctx).
		Joins("JOIN teacher_classes ON teacher_classes.teacher_id = users.id").
		Where("teacher_classes.class_id = ?", classID).
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("user repository: find teachers by class id: %w", err)
	}
	users := make([]*entity.User, len(ms))
	for i, m := range ms {
		m := m
		users[i] = m.ToEntity()
	}
	return users, nil
}

func (r *userRepository) IsTeacherOfClass(ctx context.Context, teacherID, classID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("teacher_classes").
		Where("teacher_id = ? AND class_id = ?", teacherID, classID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("user repository: is teacher of class: %w", err)
	}
	return count > 0, nil
}
