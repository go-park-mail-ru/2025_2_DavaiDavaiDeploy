package models

type SignUpInput struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar,omitempty"`
	Country     string `json:"country,omitempty"`
}
