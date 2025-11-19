package models

import "html"

type EnableTwoFactorResponse struct {
	Has2FA bool   `json:"has_2fa" binding:"required"`
	QrCode string `json:"qr_code" binding:"required"`
}

func (etfr *EnableTwoFactorResponse) Sanitize() {
	etfr.QrCode = html.EscapeString(etfr.QrCode)
}
