package presistence

type (
	Role string
)

const (
	StaffRole Role = "staff"
	AdminRole Role = "admin"
	UserRole  Role = "user"
	AllRole   Role = "all"
)
