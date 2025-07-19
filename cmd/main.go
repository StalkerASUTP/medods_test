package main

import (
	"authentication-app/internal/config"
	"authentication-app/internal/handler/deactivate"
	"authentication-app/internal/handler/getid"
	"authentication-app/internal/handler/save"
	"authentication-app/internal/repository"
	"authentication-app/internal/util/tokens"
	"authentication-app/storage/db"

	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	logger.Info("logger initialized")
	conf := config.NewConfig()
	logger.Info("config loaded", "config", conf)
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DBName) //postgres://<username>:<password>@<host>:<port>/<dbname>

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
	}
	defer conn.Close(ctx)
	db := db.New(conn)
	logger.Info("database connected")
	repo := repository.NewRepository(db, conn)
	logger.Info("repository initialized")

	logger.Info("config loaded")
	tokenManager, err := tokens.NewManager(conf.SecretKey)
	if err != nil {
		logger.Error("failed to create token manager", "error", err)
	}
	logger.Info("token manager initialized")
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/tokens", save.New(logger, repo, conf, tokenManager)).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access token refreshed\n"))
		w.Write([]byte("Refresh token generated\n"))
	})
	authRouter.HandleFunc("/user", getid.New(logger, repo,tokenManager)).Methods(http.MethodGet)
	authRouter.HandleFunc("/logout", deactivate.New(logger,repo,tokenManager)).Methods(http.MethodPatch)
	logger.Info("router initialized")
	server := &http.Server{
		Addr:         conf.Address,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
		IdleTimeout:  conf.IdleTimeout,
		Handler:      router,
	}
	logger.Info("server initialized")
	logger.Info("server is listening", "address", conf.HTTPServer.Address)
	if err = server.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "error", err)
	}

}
