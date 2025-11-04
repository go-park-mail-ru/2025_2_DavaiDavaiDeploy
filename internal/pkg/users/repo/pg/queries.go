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
