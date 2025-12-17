package models

type Pager struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

func NewPager(count, offset int) Pager {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	return Pager{
		Count:  count,
		Offset: offset,
	}
}
