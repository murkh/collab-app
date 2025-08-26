package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	jwt "github.com/murkh/collab-app/collab-app/internal/auth"
	"github.com/murkh/collab-app/collab-app/internal/handler"
	"github.com/murkh/collab-app/collab-app/internal/store"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN required")
	}

	privKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privKeyPath == "" {
		log.Fatal("JWT_PRIVATE_KEY_PATH required")
	}

	db, err := store.NewPostgresStore(context.Background(), dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close(context.Background())

	signer, err := jwt.NewSignerFromFile(privKeyPath)
	if err != nil {
		log.Fatalf("load signer: %v", err)
	}

	h := handler.NewHandler(db, signer)

	r := chi.NewRouter()
	r.Post("/api/collab/token", h.IssueCollabToken)

	addr := ":8080"
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("API listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

}
