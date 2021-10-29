package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/0xdod/trove"
)

type MockUserService struct {
	Users      []*trove.User
	CreateUser func(*trove.User) error
}

func (m *MockUserService) Create(_ context.Context, u *trove.User) error {
	return m.CreateUser(u)
}

func (m *MockUserService) FindUserByID(_ context.Context, _ int) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m *MockUserService) FindUsers(_ context.Context, _ trove.UserFilter) ([]*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m *MockUserService) UpdateUser(_ context.Context, _ int, _ trove.UserPatch) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (m *MockUserService) DeleteUser(_ context.Context, _ int) error {
	panic("not implemented") // TODO: Implement
}

func TestSignup(t *testing.T) {
	userJsonReq := `
{
"first_name": "test",
"last_name": "user",
"email": "test@user.com",
"password": "123456"
}
	`
	sr := strings.NewReader(userJsonReq)
	req, _ := http.NewRequest(http.MethodPost, "/users", sr)
	w := httptest.NewRecorder()
	mus := &MockUserService{}

	mus.CreateUser = func(u *trove.User) error {
		u.ID = 1
		mus.Users = append(mus.Users, u)
		return nil
	}

	srv := &Server{UserService: mus}

	srv.RegisterUser(w, req)

	expectedResp := &trove.User{
		ID:        1,
		FirstName: "test",
		LastName:  "user",
		Email:     "test@user.com",
	}

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	assertStringsEqual(t, mus.Users[0].Email, expectedResp.Email)
	gotResp := new(trove.User)
	_ = srv.readJSON(w.Body, gotResp)
	assertDeepEqual(t, gotResp, expectedResp)
}

func assertDeepEqual(tb testing.TB, got, want interface{}) {
	tb.Helper()

	if !reflect.DeepEqual(got, want) {
		tb.Errorf("expected %#+v, but got %#+v", want, got)
	}

}

func assertStringsEqual(tb testing.TB, got, want string) {
	tb.Helper()

	if got != want {
		tb.Errorf("expected %q, got %q", want, got)
	}
}
