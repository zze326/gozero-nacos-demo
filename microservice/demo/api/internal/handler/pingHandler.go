package handler

import (
	"demo/common/response"
	"net/http"

	"demo/microservice/demo/api/internal/logic"
	"demo/microservice/demo/api/internal/svc"
)

func pingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewPingLogic(r.Context(), ctx)
		err := l.Ping()
		response.Response(w, nil, err, r.Context())
	}
}
