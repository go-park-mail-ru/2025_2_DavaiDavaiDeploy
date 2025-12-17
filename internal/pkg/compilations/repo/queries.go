package repo

import _ "embed"

//go:embed sql/getCompilationByIDQuery.sql
var GetCompilationByIDQuery string

//go:embed sql/getCompilationsWithPaginationQuery.sql
var GetCompilationsWithPaginationQuery string

//go:embed sql/getFilmsByCompilationQuery.sql
var GetFilmsByCompilationQuery string
