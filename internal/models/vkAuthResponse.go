package models

import "html"

type VKAuthResponse struct {
	User struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
		Sex       int    `json:"sex"`
		Verified  bool   `json:"verified"`
		Birthday  string `json:"birthday"`
	} `json:"user"`
}

func (u *VKAuthResponse) Sanitize() {
	u.User.UserID = html.EscapeString(u.User.UserID)
	u.User.FirstName = html.EscapeString(u.User.FirstName)
	u.User.LastName = html.EscapeString(u.User.LastName)
	u.User.Avatar = html.EscapeString(u.User.Avatar)
	u.User.Birthday = html.EscapeString(u.User.Birthday)
}
