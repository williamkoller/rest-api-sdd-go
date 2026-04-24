package dto

import "github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"

type ClassroomResponse struct {
	ID       string `json:"id"`
	UnitID   string `json:"unitId"`
	Code     string `json:"code"`
	Capacity int    `json:"capacity"`
	Active   bool   `json:"active"`
}

func MapClassroom(e *entity.Classroom) *ClassroomResponse {
	return &ClassroomResponse{
		ID:       e.ID,
		UnitID:   e.UnitID,
		Code:     e.Code,
		Capacity: e.Capacity,
		Active:   e.Active,
	}
}

func MapClassrooms(list []*entity.Classroom) []*ClassroomResponse {
	res := make([]*ClassroomResponse, len(list))
	for i, e := range list {
		res[i] = MapClassroom(e)
	}
	return res
}
