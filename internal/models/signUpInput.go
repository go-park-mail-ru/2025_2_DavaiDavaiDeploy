package models

import "html"

type SignUpInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (sui *SignUpInput) Sanitize() {
	sui.Login = html.EscapeString(sui.Login)
	sui.Password = html.EscapeString(sui.Password)
}
