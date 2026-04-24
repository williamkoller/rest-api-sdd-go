package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type StudentResponse struct {
	ID                 string    `json:"id"`
	SchoolID           string    `json:"schoolId"`
	Name               string    `json:"name"`
	BirthDate          time.Time `json:"birthDate"`
	CPF                string    `json:"cpf"`
	RegistrationNumber string    `json:"registrationNumber"`
	Active             bool      `json:"active"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type EnrollmentResponse struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"studentId"`
	ClassID      string    `json:"classId"`
	AcademicYear int       `json:"academicYear"`
	EnrolledAt   time.Time `json:"enrolledAt"`
	Status       string    `json:"status"`
}

func MapStudent(e *entity.Student) *StudentResponse {
	return &StudentResponse{
		ID:                 e.ID,
		SchoolID:           e.SchoolID,
		Name:               e.Name,
		BirthDate:          e.BirthDate,
		CPF:                formatCPF(e.CPF),
		RegistrationNumber: e.RegistrationNumber,
		Active:             e.Active,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func MapStudents(list []*entity.Student) []*StudentResponse {
	res := make([]*StudentResponse, len(list))
	for i, e := range list {
		res[i] = MapStudent(e)
	}
	return res
}

func MapEnrollment(e *entity.Enrollment) *EnrollmentResponse {
	return &EnrollmentResponse{
		ID:           e.ID,
		StudentID:    e.StudentID,
		ClassID:      e.ClassID,
		AcademicYear: e.AcademicYear,
		EnrolledAt:   e.EnrolledAt,
		Status:       string(e.Status),
	}
}
