package repo

import _ "embed"

//go:embed sql/getUserByIDQuery.sql
var GetUserByIDQuery string

//go:embed sql/getUserByLoginQuery.sql
var GetUserByLoginQuery string

//go:embed sql/updateUserPasswordQuery.sql
var UpdateUserPasswordQuery string

//go:embed sql/updateUserAvatarQuery.sql
var UpdateUserAvatarQuery string

//go:embed sql/createFeedbackQuery.sql
var CreateFeedbackQuery string

//go:embed sql/getFeedbackByIDQuery.sql
var GetFeedbackByIDQuery string

//go:embed sql/getFeedbacksByUserIDQuery.sql
var GetFeedbacksByUserIDQuery string

//go:embed sql/updateFeedbackQuery.sql
var UpdateFeedbackQuery string

//go:embed sql/getFeedbackStatsQuery.sql
var GetFeedbackStatsQuery string
