package repo

import _ "embed"

//go:embed sql/getActorByIDQuery.sql
var GetActorByID string

//go:embed sql/getActorFilmsCountQuery.sql
var GetActorFilmsCount string

//go:embed sql/getFilmAvgRatingQuery.sql
var GetFilmAvgRating string

//go:embed sql/getFilmsByActorQuery.sql
var GetFilmsByActor string
