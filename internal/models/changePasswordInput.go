package models

import "html"

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (cpi *ChangePasswordInput) Sanitize() {
	cpi.OldPassword = html.EscapeString(cpi.OldPassword)
	cpi.NewPassword = html.EscapeString(cpi.NewPassword)
}
