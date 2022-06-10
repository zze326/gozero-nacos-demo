package response

/**
 * @Author: zze
 * @Date: 2022/5/19 15:54
 * @Desc:
 */

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error, ctx context.Context) {
	var body Body
	if err != nil {
		body.Code = 0
		body.Msg = err.Error()
	} else {
		body.Code = 1
		body.Msg = "ok"
		body.Data = resp
	}
	body.TraceID = traceIDFromContext(ctx)
	httpx.OkJson(w, body)
}

func traceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
