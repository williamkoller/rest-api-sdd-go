package entity

import "time"

type Student struct {
	ID                 string
	SchoolID           string
	Name               string
	BirthDate          time.Time
	CPF                string
	RegistrationNumber string
	Active             bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
