package models

type SignInInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
