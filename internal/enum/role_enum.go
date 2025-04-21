package enum

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

func IsValidRole(role Role) bool {
	switch role {
	case RoleEmployee, RoleModerator:
		return true
	default:
		return false
	}
}

func (r Role) String() string {
	return string(r)
}
