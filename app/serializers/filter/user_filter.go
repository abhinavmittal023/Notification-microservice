package filter

// User struct is the serializer for user filter
type User struct {
	ID        uint   `form:"user_id"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Email     string `form:"email"`
	Verified  int    `form:"verified"`
	Role      int    `form:"role"`
}
