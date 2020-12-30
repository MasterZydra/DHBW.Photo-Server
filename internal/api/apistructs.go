package api

type RegisterReq struct {
	Username             string
	Password             string
	PasswordConfirmation string
}

type RegisterRes struct {
	Username string
	Error    string
}
