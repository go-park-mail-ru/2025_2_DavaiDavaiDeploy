package repo

import _ "embed"

//go:embed sql/getGenreByIDQuery.sql
var GetGenreByIDQuery string

//go:embed sql/getGenresWithPaginationQuery.sql
var GetGenresWithPaginationQuery string

//go:embed sql/getFilmAvgRatingQuery.sql
var GetFilmAvgRatingQuery string

//go:embed sql/getFilmsByGenreQuery.sql
var GetFilmsByGenreQuery string
