package sessions

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/mtlynch/jeff"
	jeff_sqlite "github.com/mtlynch/jeff/sqlite"
)

type (
	ManagerParams struct {
		RequireTLS bool
	}

	Option func(*ManagerParams)

	Manager struct {
		j *jeff.Jeff
	}
)

var defaultParams = ManagerParams{
	RequireTLS: true,
}

func NewManager(sqliteDB *sql.DB, options ...Option) (Manager, error) {
	store, err := jeff_sqlite.New(sqliteDB)
	if err != nil {
		return Manager{}, err
	}
	params := defaultParams
	for _, o := range options {
		o(&params)
	}
	jeff_opts := []func(*jeff.Jeff){jeff.CookieName("token")}
	if !params.RequireTLS {
		jeff_opts = append(jeff_opts, jeff.Insecure)
	}
	j := jeff.New(store, jeff_opts...)
	return Manager{
		j: j,
	}, nil
}

// AllowNonTlsConnections allows non-TLS connections for a session, removing the
// Secure attribute from session cookies.
func AllowNonTlsConnections() func(p *ManagerParams) {
	return func(p *ManagerParams) {
		p.RequireTLS = false
	}
}

func (m Manager) CreateSession(w http.ResponseWriter, ctx context.Context, key Key, session Session) error {
	return m.j.Set(ctx, w, key.Bytes(), session)
}

func (m Manager) SessionFromContext(ctx context.Context) (Session, error) {
	sess := jeff.ActiveSession(ctx)
	if len(sess.Key) == 0 {
		return Session{}, ErrNoSessionFound
	}

	return sess.Meta, nil
}

func (m Manager) EndSession(ctx context.Context, w http.ResponseWriter) {
	sess := jeff.ActiveSession(ctx)
	if len(sess.Key) > 0 {
		if err := m.j.Delete(ctx, sess.Key); err != nil {
			log.Printf("failed to delete session: %v", err)
		}
	}

	if err := m.j.Clear(ctx, w); err != nil {
		log.Printf("failed to clear session: %v", err)
	}
}

func (m Manager) WrapRequest(next http.Handler) http.Handler {
	return m.j.Public(next)
}
