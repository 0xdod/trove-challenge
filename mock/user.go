package mock

import (
	"context"

	"github.com/0xdod/trove"
)

type UserService struct {
	Users             []*trove.User
	CreateUserFn      func(*trove.User) error
	UpdateUserFn      func(int, trove.UserPatch) (*trove.User, error)
	FindUserByTokenFn func() *trove.User
}

func (m UserService) Create(_ context.Context, u *trove.User) error {
	return m.CreateUserFn(u)
}

func (m UserService) FindUserByID(_ context.Context, _ int) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m UserService) FindUserByToken(_ context.Context, _ string) (*trove.User, error) {
	return m.FindUserByTokenFn(), nil
}

func (m UserService) FindUserByEmail(_ context.Context, _ string) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m UserService) FindUsers(_ context.Context, _ trove.UserFilter) ([]*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m UserService) UpdateUser(_ context.Context, id int, up trove.UserPatch) (*trove.User, error) {
	return m.UpdateUserFn(id, up)
}

func (m UserService) DeleteUser(_ context.Context, _ int) error {
	panic("not implemented") // TODO: Implement
}
