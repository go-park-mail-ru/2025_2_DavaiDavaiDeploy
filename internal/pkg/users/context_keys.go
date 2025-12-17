package users

type contextKey string

func (c contextKey) String() string {
	return "auth context key " + string(c)
}

const (
	UserKey contextKey = "user_id"
)
