package models

import "html"

type SignInInput struct {
	Login    string  `json:"login" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Code     *string `json:"qr_code,omitempty" `
}

func (sii *SignInInput) Sanitize() {
	sii.Login = html.EscapeString(sii.Login)
	sii.Password = html.EscapeString(sii.Password)
	if sii.Code != nil {
		sanitized := html.EscapeString(*sii.Code)
		sii.Code = &sanitized
	}
}
