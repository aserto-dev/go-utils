package grpcutil

type CtxKey string

// If a call is authenticated with a machine account, these hold info about the connection and cert that resolved
// that machine account.
var (
	ValueMachineAccountConnectionID = CtxKey("machine-account-connection-id")
	ValueMachineAccountCertID       = CtxKey("machine-account-cert-id")
)
