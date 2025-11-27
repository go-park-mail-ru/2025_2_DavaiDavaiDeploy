package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CursorPager struct {
	CreatedAt *timestamppb.Timestamp `json:"created_at"`
	Count     int                    `json:"count"`
}

func NewCursorPager(created_at *timestamppb.Timestamp, count int) CursorPager {
	if count < 0 {
		count = 0
	}
	return CursorPager{
		Count:     count,
		CreatedAt: created_at,
	}
}
