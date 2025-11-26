package models

type DisableTwoFactorResponse struct {
	Has2FA bool `json:"has_2fa" binding:"required"`
}
