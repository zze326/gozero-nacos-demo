package handler

import (
	"demo/common/response"
	"demo/microservice/demo/api/internal/logic"
	"demo/microservice/demo/api/internal/svc"
	"demo/microservice/demo/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func getUserHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewGetUserLogic(r.Context(), ctx)
		resp, err := l.GetUser(&req)
		response.Response(w, resp, err, r.Context())
	}
}
