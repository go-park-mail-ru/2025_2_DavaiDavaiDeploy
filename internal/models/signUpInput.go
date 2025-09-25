package models

type SignUpInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
