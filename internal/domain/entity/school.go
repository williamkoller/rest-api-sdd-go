package entity

import "time"

type School struct {
	ID        string
	Name      string
	CNPJ      string
	Email     string
	Phone     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
