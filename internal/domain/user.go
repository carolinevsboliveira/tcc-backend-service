package domain

type Role string

const (
	RoleStudent   Role = "student"
	RoleProfessor Role = "professor"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	DepartmentID string
	Active       bool
}
