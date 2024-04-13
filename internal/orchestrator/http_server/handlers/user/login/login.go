package login

import (
	"context"
	"log/slog"
	"net/http"
)

func New(logger *slog.Logger, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
