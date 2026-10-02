// Comando server inicia a API REST de tarefas.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaoandradecuiaba/ccb/backend/internal/httpapi"
	"github.com/joaoandradecuiaba/ccb/backend/internal/task"
)

func main() {
	addr := getenv("ADDR", ":8080")
	dataFile := getenv("DATA_FILE", "data/tasks.json")

	store, err := task.NewStore(dataFile)
	if err != nil {
		log.Fatalf("falha ao carregar dados: %v", err)
	}

	srv := &http.Server{
		Addr: addr,
		Handler: httpapi.NewHandler(store, httpapi.Config{
			AllowedOrigin: getenv("CORS_ORIGIN", "http://localhost:4200"),
			StaticDir:     os.Getenv("STATIC_DIR"),
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("servidor ouvindo em %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("encerrando...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro ao encerrar: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
