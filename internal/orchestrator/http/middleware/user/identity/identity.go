package identity

import (
	"context"
	"github.com/go-chi/render"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/lib/jwt"
	"net/http"
	"strings"
)

func New(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				render.JSON(w, r, resp.Error("Authorization header is empty"))

				return
			}
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				render.JSON(w, r, resp.Error("Authorization header is invalid"))
				return
			}

			if len(headerParts[1]) == 0 {
				render.JSON(w, r, resp.Error("Authorization token is empty"))
				return
			}

			userId, err := jwt.GetUserID(headerParts[1], secret)
			if err != nil {
				render.JSON(w, r, resp.Error(err.Error()))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func GetUserID(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}
