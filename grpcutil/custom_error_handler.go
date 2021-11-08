package grpcutil

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	HttpStatusErrorMetadata = "aserto-http-statuscode"
)

func CustomErrorHandler(ctx context.Context, gtw *runtime.ServeMux, ms runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		st := status.Convert(err)
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.ErrorInfo:
				value, ok := t.Metadata[HttpStatusErrorMetadata]
				if ok {
					code, conv_err := strconv.Atoi(value)
					if conv_err != nil {
						logger := zerolog.Ctx(ctx)
						logger.Error().Err(conv_err).Msg("Failed to detect http status code associated with this AsertoErrror")
					} else {
						var httpStatusError runtime.HTTPStatusError
						httpStatusError.Err = err
						httpStatusError.HTTPStatus = code
						runtime.DefaultHTTPErrorHandler(ctx, gtw, ms, w, r, &httpStatusError)
						return
					}
				}
			}
		}
	}

	runtime.DefaultHTTPErrorHandler(ctx, gtw, ms, w, r, err)
}
