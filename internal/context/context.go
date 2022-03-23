package context

import (
	"context"
	repo "github.com/3n0ugh/kalenderium/pkg/account/repository"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

// SetUser method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func SetUser(r *http.Request, u *repo.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, u)
	return r.WithContext(ctx)
}

// GetUser retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
func GetUser(r *http.Request) *repo.User {
	user, ok := r.Context().Value(userContextKey).(*repo.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
