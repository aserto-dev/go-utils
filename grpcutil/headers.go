package grpcutil

import (
	"context"
	"net/textproto"
	"strings"
)

var (
	HeaderAsertoTenantID  = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Tenant-Id"))
	HeaderAsertoAccountID = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Account-Id"))
	HeaderAsertoTenantKey = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Api-Key"))
	HeaderAsertoRequestID = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Request-Id"))

	HeaderAsertoTenantIDLowercase  = CtxKey(strings.ToLower(string(HeaderAsertoTenantID)))
	HeaderAsertoAccountIDLowercase = CtxKey(strings.ToLower(string(HeaderAsertoAccountID)))
	HeaderAsertoTenantKeyLowercase = CtxKey(strings.ToLower(string(HeaderAsertoTenantKey)))
	HeaderAsertoRequestIDLowercase = CtxKey(strings.ToLower(string(HeaderAsertoRequestID)))
)

func ExtractTenantID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoTenantID).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractAccountID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoAccountID).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractTenantKey(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoTenantKey).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractRequestID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoRequestID).(string)
	if !ok {
		return ""
	}

	return id
}
