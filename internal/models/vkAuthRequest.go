package models

import "html"

type VKAuthRequest struct {
	Login       string `json:"login"`
	AccessToken string `json:"access_token"`
}

func (u *VKAuthRequest) Sanitize() {
	u.Login = html.EscapeString(u.Login)
	u.AccessToken = html.EscapeString(u.AccessToken)
}
