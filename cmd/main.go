package main

import (
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"time"

	"surfe-actions/internal/actions"
	actionsrepository "surfe-actions/internal/actions/repository"
	"surfe-actions/internal/config"
	http2 "surfe-actions/internal/http"
	"surfe-actions/internal/users"
	usersrepository "surfe-actions/internal/users/repository"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		stdlog.Fatalf("unable to process config: %v", err)
	}

	usersRepo, err := usersrepository.NewRepository()
	if err != nil {
		stdlog.Fatalf("unable to process users: %v", err)
	}

	actionsRepo, err := actionsrepository.NewRepository()
	if err != nil {
		stdlog.Fatalf("unable to process actions: %v", err)
	}

	usersService := users.NewService(usersRepo)
	actionsService := actions.NewService(actionsRepo)

	mux := http2.NewRouter(usersService, actionsService)

	err = NewServer(mux, fmt.Sprintf(":%d", cfg.Rest.Port)).ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		stdlog.Println(err)
	}
}

func NewServer(router http.Handler, addr string) *http.Server {
	const readHeaderTimeout = time.Second * 2

	return &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}
