package grpcutil

import (
	"context"
	"errors"
	"net/http"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func CustomErrorHandler(ctx context.Context, gtw *runtime.ServeMux, ms runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	var httpStatusError *runtime.HTTPStatusError
	var custom *cerr.AsertoError

	if errors.As(err, &custom) {
		httpStatusError.Err = err
		httpStatusError.HTTPStatus = custom.HttpCode
	} else {
		httpStatusError.Err = err
	}
	runtime.DefaultHTTPErrorHandler(ctx, gtw, ms, w, r, httpStatusError)
}
