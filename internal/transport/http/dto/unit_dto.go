package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type UnitResponse struct {
	ID        string    `json:"id"`
	SchoolID  string    `json:"schoolId"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zipCode"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func MapUnit(e *entity.Unit) *UnitResponse {
	return &UnitResponse{
		ID:        e.ID,
		SchoolID:  e.SchoolID,
		Name:      e.Name,
		Address:   e.Address,
		City:      e.City,
		State:     e.State,
		ZipCode:   formatZipCode(e.ZipCode),
		Phone:     formatPhone(e.Phone),
		Active:    e.Active,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func MapUnits(list []*entity.Unit) []*UnitResponse {
	res := make([]*UnitResponse, len(list))
	for i, e := range list {
		res[i] = MapUnit(e)
	}
	return res
}
