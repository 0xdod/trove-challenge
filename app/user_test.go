package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/0xdod/trove"
	"github.com/0xdod/trove/mock"
	"github.com/gorilla/mux"
)

func Test_registerUser(t *testing.T) {
	t.Run("test successful user creation", func(t *testing.T) {

		userJsonReq := `
{
"first_name": "test",
"last_name": "user",
"email": "test@user.com",
"password": "12345678"
}
	`
		sr := strings.NewReader(userJsonReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", sr)
		w := httptest.NewRecorder()
		mus := &mock.UserService{}

		mus.CreateUserFn = func(u *trove.User) error {
			u.ID = 1
			mus.Users = append(mus.Users, u)
			return nil
		}

		expectedResp := trove.User{
			ID:        1,
			FirstName: "test",
			LastName:  "user",
			Email:     "test@user.com",
		}

		srv := &Server{
			router:      mux.NewRouter().StrictSlash(true),
			UserService: mus,
		}

		srv.routes()

		srv.router.ServeHTTP(w, req)

		gotResp := trove.User{}

		extractDataFromResponse(w.Body, &gotResp)

		_ = srv.readJSON(w.Body, gotResp)

		if ct := w.Header().Get("content-type"); ct != "application/json" {
			t.Errorf("expected content-type %q, but got %q", "application/json", ct)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		assertDeepEqual(t, gotResp, expectedResp)
	})
}

func Test_updateUser(t *testing.T) {
	t.Run("test successful user update", func(t *testing.T) {
		userJsonReq := `
{
"first_name": "changed first_name",
"last_name": "changed last_name"
}
`
		sr := strings.NewReader(userJsonReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users/1", sr)
		req.Header.Add("Authorization", "Bearer HelloWorld")
		w := httptest.NewRecorder()
		mus := &mock.UserService{
			Users: []*trove.User{
				{
					ID:        1,
					FirstName: "test",
					LastName:  "user",
					Email:     "test@user.com",
				},
			},
		}

		mus.UpdateUserFn = func(id int, p trove.UserPatch) (*trove.User, error) {
			var user *trove.User

			for _, u := range mus.Users {
				if u.ID == id {
					user = u
					break
				}
			}
			user.ApplyPatch(p)
			return user, nil
		}

		mus.FindUserByTokenFn = func() *trove.User {
			return mus.Users[0]
		}

		expectedResp := trove.User{
			ID:        1,
			FirstName: "changed first_name",
			LastName:  "changed last_name",
			Email:     "test@user.com",
		}

		srv := &Server{
			router:      mux.NewRouter().StrictSlash(true),
			UserService: mus,
		}

		srv.routes()

		srv.router.ServeHTTP(w, req)

		gotResp := trove.User{}

		extractDataFromResponse(w.Body, &gotResp)

		_ = srv.readJSON(w.Body, gotResp)

		if ct := w.Header().Get("content-type"); ct != "application/json" {
			t.Errorf("expected content-type %q, but got %q", "application/json", ct)
		}

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		assertDeepEqual(t, gotResp, expectedResp)
	})
}

func assertDeepEqual(tb testing.TB, got, want interface{}) {
	tb.Helper()

	if !reflect.DeepEqual(got, want) {
		tb.Errorf("expected %#+v, but got %#+v", want, got)
	}

}

// func assertStringsEqual(tb testing.TB, got, want string) {
// 	tb.Helper()

// 	if got != want {
// 		tb.Errorf("expected %q, got %q", want, got)
// 	}
// }

func extractDataFromResponse(body io.Reader, v interface{}) error {
	grm := genericResponseModel{}
	_ = json.NewDecoder(body).Decode(&grm)
	tempBuf := new(bytes.Buffer)
	_ = json.NewEncoder(tempBuf).Encode(grm.Data)
	_ = json.NewDecoder(tempBuf).Decode(v)
	return nil
}
