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

type ImageReq struct {
	Data string
}

type ImageRes struct {
	Data string
}
