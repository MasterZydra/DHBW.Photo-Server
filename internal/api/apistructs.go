package api

import "net/http"

type RegisterReq struct {
	Username             string
	Password             string
	PasswordConfirmation string
}

type RegisterRes struct {
	Username string
	Error    string
}

type LoginReq struct {
	Username string
	Password string
}

type LoginRes struct {
	Cookie http.Cookie
	Error  string
}
