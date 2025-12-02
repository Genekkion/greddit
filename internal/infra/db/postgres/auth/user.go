package authdb

import (
	"context"
	"greddit/internal/domains/auth"
	"greddit/internal/infra/db/postgres"
)

// UsersRepo implements the forumports.UsersRepo interface.
type UsersRepo struct {
	postgres.BaseRepo
}

func (u UsersRepo) CreateUser(ctx context.Context, value auth.UserValue) (user *auth.User, err error) {

	return nil, nil
}

func (u UsersRepo) GetUserById(ctx context.Context, id auth.UserId) (user *auth.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UsersRepo) UpdateDisplayName(ctx context.Context, id auth.UserId, displayName string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u UsersRepo) DeleteUser(ctx context.Context, id auth.UserId) (err error) {
	//TODO implement me
	panic("implement me")
}
