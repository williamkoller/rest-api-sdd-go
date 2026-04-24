package entity

import "time"

type Role string

const (
	RoleSuperAdmin  Role = "super_admin"
	RoleSchoolAdmin Role = "school_admin"
	RoleUnitStaff   Role = "unit_staff"
	RoleTeacher     Role = "teacher"
	RoleGuardian    Role = "guardian"
)

type User struct {
	ID           string
	SchoolID     string
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
