package repo

import _ "embed"

//go:embed sql/getFilmsFromSearchQuery.sql
var GetFilmsFromSearchQuery string

//go:embed sql/getActorsFromSearchQuery.sql
var GetActorsFromSearchQuery string
