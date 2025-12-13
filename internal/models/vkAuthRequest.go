package models

import "html"

type VKAuthRequest struct {
	Login       *string `json:"login"`
	AccessToken string  `json:"access_token"`
}

func (u *VKAuthRequest) Sanitize() {
	if u.Login != nil {
		sanitized := html.EscapeString(*u.Login)
		u.Login = &sanitized
	}
	u.AccessToken = html.EscapeString(u.AccessToken)
}
