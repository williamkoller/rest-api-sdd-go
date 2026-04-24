package entity

import "time"

type Unit struct {
	ID        string
	SchoolID  string
	Name      string
	Address   string
	City      string
	State     string
	ZipCode   string
	Phone     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
