package api

import (
	auth "auditlog/internal/authentication"
	"auditlog/internal/config"
	"auditlog/internal/data"
	"auditlog/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type App struct {
	L    *utils.AuditLogLogger
	R    *httprouter.Router
	DB   data.DBAccessor
	Auth auth.Authenticator
}

// factory function for new App instance
func New(l *utils.AuditLogLogger, db data.DBAccessor, auth auth.Authenticator) App {
	app := App{
		L:    l,
		DB:   db,
		Auth: auth,
	}
	r := app.registerRoutes()
	app.R = r

	return app
}

// Establishes an HTTP server listening on the configured port using
// the application's router for request routing
func (a *App) Run(ctx context.Context) {

	// TODO: include TLS certificates to ListenAndServeTLS
	// omitted for convenience
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.AppPort),
		Handler:      a.R,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			a.L.Err(fmt.Sprintf("%v", err))
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan // blocks
	a.L.Info(fmt.Sprintf("Shutting down server %v", sig))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
