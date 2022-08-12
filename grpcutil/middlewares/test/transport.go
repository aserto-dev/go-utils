package test

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TestServerTransportStream struct {
	method string
}

func ServerTransportStream(method string) *TestServerTransportStream {
	return &TestServerTransportStream{method}
}

func (ts *TestServerTransportStream) Method() string {
	return ts.method
}
func (ts *TestServerTransportStream) SetHeader(md metadata.MD) error {
	return nil
}
func (ts *TestServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}
func (ts *TestServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}

type TestServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func ServerStream(ctx context.Context) *TestServerStream {
	return &TestServerStream{ctx: ctx}
}

func (ss *TestServerStream) Context() context.Context {
	return ss.ctx
}

func (ss *TestServerStream) SendMsg(m interface{}) error {
	return nil
}

func (f *TestServerStream) RecvMsg(m interface{}) error {
	return nil
}

func (f *TestServerStream) Method() string {
	return ""
}
