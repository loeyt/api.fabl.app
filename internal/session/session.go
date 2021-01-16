package session

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func init() {
	gob.Register(uuid.Nil)
}

type contextKey int

const (
	ckSession contextKey = iota
)

var (
	Unauthenticated = codes.Unauthenticated
)

type middleware struct {
	*runtime.ServeMux
	*sessions.CookieStore
	name string
}

// Wrap a runtime.ServeMux to use Cookie based authentication. Should be
// combined with ForwardResponseOption.
func Wrap(mux *runtime.ServeMux, name string, keyPairs ...[]byte) http.Handler {
	return &middleware{
		ServeMux:    mux,
		CookieStore: sessions.NewCookieStore(keyPairs...),
		name:        "session",
	}
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Disable session handling if the Authorization handler is being used.
	if r.Header.Get("Authorization") != "" {
		m.ServeMux.ServeHTTP(w, r)
		return
	}
	session, err := m.CookieStore.Get(r, m.name)
	if err != nil {
		// TODO: log it?
		m.ServeMux.ServeHTTP(w, r)
		return
	}
	session.Options.HttpOnly = true
	m.ServeMux.ServeHTTP(w, r.WithContext(
		context.WithValue(r.Context(), ckSession, session)),
	)
}

func Account(ctx context.Context) (uuid.UUID, error) {
	session, ok := ctx.Value(ckSession).(*sessions.Session)
	if !ok {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "no session")
	}
	if id, ok := session.Values["account"].(uuid.UUID); ok {
		return id, nil
	}
	return uuid.Nil, status.Errorf(codes.Unauthenticated, "no account in session")
}

func Login(ctx context.Context, id uuid.UUID) error {
	session, ok := ctx.Value(ckSession).(*sessions.Session)
	if !ok {
		// TODO: invocation or configuration error: no session in ctx
		return nil
	}
	session.Values["account"] = id
	return nil
}

func Logout(ctx context.Context) error {
	session, ok := ctx.Value(ckSession).(*sessions.Session)
	if !ok {
		// TODO: invocation or configuration error: no session in ctx
		return nil
	}
	delete(session.Values, "account")
	return nil
}

func ForwardResponseOption(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	session, ok := ctx.Value(ckSession).(*sessions.Session)
	if !ok {
		return nil
	}
	if _, ok := session.Values["account"]; session.IsNew && !ok {
		// Don't save a session if we entered unauthenticated, and leave the
		// same.
		return nil
	}
	err := session.Save(nil, w)
	if err != nil {
		// TODO: log? return? I don't know...
		return err
	}
	return nil
}
