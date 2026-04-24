package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type SchoolResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func MapSchool(e *entity.School) *SchoolResponse {
	return &SchoolResponse{
		ID:        e.ID,
		Name:      e.Name,
		CNPJ:      formatCNPJ(e.CNPJ),
		Email:     e.Email,
		Phone:     formatPhone(e.Phone),
		Active:    e.Active,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func MapSchools(list []*entity.School) []*SchoolResponse {
	res := make([]*SchoolResponse, len(list))
	for i, e := range list {
		res[i] = MapSchool(e)
	}
	return res
}
