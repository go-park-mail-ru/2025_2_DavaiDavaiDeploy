package repo

import _ "embed"

//go:embed sql/checkUserExistsQuery.sql
var CheckUserExistsQuery string

//go:embed sql/createUserQuery.sql
var CreateUserQuery string

//go:embed sql/checkUserLoginQuery.sql
var CheckUserLoginQuery string

//go:embed sql/incrementUserVersionQuery.sql
var IncrementUserVersionQuery string

//go:embed sql/getUserByLoginQuery.sql
var GetUserByLoginQuery string
