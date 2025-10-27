package models

import "html"

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (sii *SignInInput) Sanitize() {
	sii.Login = html.EscapeString(sii.Login)
	sii.Password = html.EscapeString(sii.Password)
}
