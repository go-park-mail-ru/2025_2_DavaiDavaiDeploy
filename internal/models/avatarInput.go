package models

import "mime/multipart"

type AvatarInput struct {
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
}
