package user

type User struct {
	Id       int
	Email    string
	FullName string
	Phone    string
	Image    string
	Role     string
	Password string
}

type Auth struct {
	Us  User
	Err bool
}
