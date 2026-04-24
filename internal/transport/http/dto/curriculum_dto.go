package dto

import "github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"

type CurriculumEntryResponse struct {
	ID        string `json:"id"`
	ClassID   string `json:"classId"`
	Subject   string `json:"subject"`
	TeacherID string `json:"teacherId"`
	DayOfWeek string `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func MapCurriculumEntry(e *entity.CurriculumEntry) *CurriculumEntryResponse {
	return &CurriculumEntryResponse{
		ID:        e.ID,
		ClassID:   e.ClassID,
		Subject:   e.Subject,
		TeacherID: e.TeacherID,
		DayOfWeek: string(e.DayOfWeek),
		StartTime: e.StartTime,
		EndTime:   e.EndTime,
	}
}

func MapCurriculumEntries(list []*entity.CurriculumEntry) []*CurriculumEntryResponse {
	res := make([]*CurriculumEntryResponse, len(list))
	for i, e := range list {
		res[i] = MapCurriculumEntry(e)
	}
	return res
}
