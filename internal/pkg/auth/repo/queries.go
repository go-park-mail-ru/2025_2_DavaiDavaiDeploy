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

//go:embed sql/enable2FaQuery.sql
var Enable2FaQuery string

//go:embed sql/disable2FaQuery.sql
var Disable2FaQuery string

//go:embed sql/getUserByIDQuery.sql
var GetUserByIDQuery string

//go:embed sql/checkUserTwoFactorQuery.sql
var CheckUserTwoFactorQuery string

//go:embed sql/checkUserSecretCodeQuery.sql
var CheckUserSecretCodeQuery string
