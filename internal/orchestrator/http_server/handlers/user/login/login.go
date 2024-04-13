package login

import (
	"context"
	"github.com/go-chi/render"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
	JWTToken string `json:"jwt_token"`
}

type Auth interface {
	Login(ctx context.Context, login, password string) (string, error)
}

func New(log *slog.Logger, ctx context.Context, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.user.login.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Login == "" {
			log.Error("login is required")

			render.JSON(w, r, resp.Error("login is required"))

			return
		}

		if req.Password == "" {
			log.Error("password is required")

			render.JSON(w, r, resp.Error("password is required"))

			return
		}

		token, err := auth.Login(ctx, req.Login, req.Password)
		if err != nil {
			log.Error("error logging in:", err)

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		responseOK(w, r, token)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		JWTToken: token,
	})
}
