package user

type User struct {
	Name     string
	Password string
}

func (u *User) ToCsv() []string {
	return []string{u.Name, u.Password}
}
