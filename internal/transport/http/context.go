package http

import (
	"context"
	"errors"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
)

type contextKey string

const userContextKey contextKey = "user"

func withUser(ctx context.Context, user domain.UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (domain.UserContext, error) {
	user, ok := ctx.Value(userContextKey).(domain.UserContext)
	if !ok {
		return domain.UserContext{}, errors.New("user not found in context")
	}

	return user, nil
}
