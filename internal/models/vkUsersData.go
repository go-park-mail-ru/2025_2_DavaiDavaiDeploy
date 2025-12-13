package models

import "html"

type VKUsersData struct {
	Login *string `json:"login"`
	VkID  string  `json:"vk_id"`
}

func (u *VKUsersData) Sanitize() {
	if u.Login != nil {
		sanitized := html.EscapeString(*u.Login)
		u.Login = &sanitized
	}
	u.VkID = html.EscapeString(u.VkID)
}
