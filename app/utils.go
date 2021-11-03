package app

import (
	"context"
	"net/http"

	"github.com/0xdod/trove"
)

func (s *Server) dummyHandlerFunc(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) coolHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (s *Server) badRequestResponse(w http.ResponseWriter, err error) {
	s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", "request error: " + err.Error(), nil})
}

func (s *Server) serverErrorResponse(w http.ResponseWriter, err error) {
	s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error: " + err.Error(), nil})
}

func UserFromContext(ctx context.Context) *trove.User {
	user, ok := ctx.Value(key("user")).(*trove.User)

	if !ok {
		return nil
	}

	return user
}
