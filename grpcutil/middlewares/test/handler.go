package test

import (
	"context"

	"google.golang.org/grpc"
)

var (
	UnaryInfo = &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}

	StreamInfo = &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsClientStream: false,
		IsServerStream: true,
	}
)

type Handler struct {
	output interface{}
	err    error
}

func NewHandler(output interface{}, err error) *Handler {
	return &Handler{output, err}
}

func (h *Handler) Unary(ctx context.Context, req interface{}) (interface{}, error) {
	return h.output, h.err
}

func (h *Handler) Stream(srv interface{}, stream grpc.ServerStream) error {
	return h.err
}
