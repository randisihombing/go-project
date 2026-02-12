package model

type Exec struct {
	ID                int
	FirstName         string
	LastName          string
	Email             string
	Username          string
	Password          string
	PasswordChangedAt string
	UserCreatedAt     string
	PasswordResetCode string
	InactiveStatus    bool
	Role              string
}
