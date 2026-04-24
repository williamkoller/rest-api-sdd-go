package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type ClassResponse struct {
	ID           string    `json:"id"`
	UnitID       string    `json:"unitId"`
	ClassroomID  string    `json:"classroomId"`
	Name         string    `json:"name"`
	GradeLevel   string    `json:"gradeLevel"`
	Shift        string    `json:"shift"`
	AcademicYear int       `json:"academicYear"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ClassDetailResponse struct {
	*ClassResponse
	EnrollmentCount int64 `json:"enrollmentCount"`
}

func MapClass(e *entity.Class) *ClassResponse {
	return &ClassResponse{
		ID:           e.ID,
		UnitID:       e.UnitID,
		ClassroomID:  e.ClassroomID,
		Name:         e.Name,
		GradeLevel:   e.GradeLevel,
		Shift:        string(e.Shift),
		AcademicYear: e.AcademicYear,
		Active:       e.Active,
		CreatedAt:    e.CreatedAt,
	}
}

func MapClassWithCount(e *usecase.ClassWithCount) *ClassDetailResponse {
	return &ClassDetailResponse{
		ClassResponse:   MapClass(e.Class),
		EnrollmentCount: e.EnrollmentCount,
	}
}

func MapClassesWithCount(list []*usecase.ClassWithCount) []*ClassDetailResponse {
	res := make([]*ClassDetailResponse, len(list))
	for i, e := range list {
		res[i] = MapClassWithCount(e)
	}
	return res
}
