package repo

import _ "embed"

//go:embed sql/getPromoFilmByIDQuery.sql
var GetPromoFilmByIDQuery string

//go:embed sql/getFilmByIDQuery.sql
var GetFilmByIDQuery string

//go:embed sql/getGenreTitleQuery.sql
var GetGenreTitleQuery string

//go:embed sql/getFilmAvgRatingQuery.sql
var GetFilmAvgRatingQuery string

//go:embed sql/getFilmsWithPaginationQuery.sql
var GetFilmsWithPaginationQuery string

//go:embed sql/getFilmPageQuery.sql
var GetFilmPageQuery string

//go:embed sql/getFilmActorsQuery.sql
var GetFilmActorsQuery string

//go:embed sql/getFilmFeedbacksQuery.sql
var GetFilmFeedbacksQuery string

//go:embed sql/checkUserFeedbackExistsQuery.sql
var CheckUserFeedbackExistsQuery string

//go:embed sql/updateFeedbackQuery.sql
var UpdateFeedbackQuery string

//go:embed sql/createFeedbackQuery.sql
var CreateFeedbackQuery string

//go:embed sql/setRatingQuery.sql
var SetRatingQuery string

//go:embed sql/getUserByLoginQuery.sql
var GetUserByLoginQuery string

//go:embed sql/insertIntoSavedQuery.sql
var InsertIntoSavedQuery string

//go:embed sql/deleteFromSavedQuery.sql
var DeleteFromSavedQuery string

//go:embed sql/checkUserLikeExistsQuery.sql
var CheckUserLikeExistsQuery string
//go:embed sql/getFilmsWithDateOfReleaseQuery.sql
var GenreetFilmsWithDateOfReleaseQuery string
