package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type WaitlistEntryResponse struct {
	ID            string    `json:"id"`
	UnitID        string    `json:"unitId"`
	GuardianName  string    `json:"guardianName"`
	GuardianEmail string    `json:"guardianEmail"`
	StudentName   string    `json:"studentName"`
	GradeLevel    string    `json:"gradeLevel"`
	AcademicYear  int       `json:"academicYear"`
	Position      int       `json:"position"`
	Status        string    `json:"status"`
	ReferralID    string    `json:"referralId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func MapWaitlistEntry(e *entity.WaitlistEntry) *WaitlistEntryResponse {
	return &WaitlistEntryResponse{
		ID:            e.ID,
		UnitID:        e.UnitID,
		GuardianName:  e.GuardianName,
		GuardianEmail: e.GuardianEmail,
		StudentName:   e.StudentName,
		GradeLevel:    e.GradeLevel,
		AcademicYear:  e.AcademicYear,
		Position:      e.Position,
		Status:        string(e.Status),
		ReferralID:    e.ReferralID,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func MapWaitlistEntries(list []*entity.WaitlistEntry) []*WaitlistEntryResponse {
	res := make([]*WaitlistEntryResponse, len(list))
	for i, e := range list {
		res[i] = MapWaitlistEntry(e)
	}
	return res
}
