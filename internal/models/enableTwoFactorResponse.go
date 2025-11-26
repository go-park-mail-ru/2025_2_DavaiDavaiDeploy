package models

type EnableTwoFactorResponse struct {
	Has2FA bool   `json:"has_2fa" binding:"required"`
	QrCode []byte `json:"qr_code" binding:"required"`
}
