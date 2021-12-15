package cerr

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MessageKey = "msg"
)

var (
	// Unknown error ID. It's returned when the implementation has not returned another AsertoError.
	ErrUnknown = newErr("E10000", codes.Internal, http.StatusInternalServerError, "an unknown error has occurred")
	// Means no tenant id was found in the current context
	ErrNoTenantID = newErr("E10001", codes.InvalidArgument, http.StatusBadRequest, "no tenant id specified")
	// Means the tenant id is not valid
	ErrInvalidTenantID = newErr("E10002", codes.InvalidArgument, http.StatusBadRequest, "invalid tenant id")
	// Means the tenant name doesn't conform to our tenant name rules
	ErrInvalidTenantName = newErr("E10003", codes.InvalidArgument, http.StatusBadRequest, "invalid tenant name")
	// Means the provider ID is invalid
	ErrInvalidProviderID = newErr("E10004", codes.InvalidArgument, http.StatusBadRequest, "invalid provider id")
	// Means the provider config name doesn't exist
	ErrInvalidProviderConfigName = newErr("E10005", codes.InvalidArgument, http.StatusBadRequest, "invalid provider config name")
	// The asked-for runtime is not yet available, but will likely be in the future.
	ErrRuntimeLoading = newErr("E10006", codes.Unavailable, http.StatusTooEarly, "runtime has not yet loaded")
	// Means a connection failed to validate.
	ErrConnectionVerification = newErr("E10007", codes.FailedPrecondition, http.StatusServiceUnavailable, "connection verification failed")
	// Returned when there's a problem retrieving a connection.
	ErrConnection = newErr("E10008", codes.Unavailable, http.StatusServiceUnavailable, "connection problem")
	// Returned when there's a problem getting a github access token.
	ErrGithubAccessToken = newErr("E10009", codes.Unavailable, http.StatusServiceUnavailable, "failed to retrieve github access token")
	// Returned when there's a problem communicating with an SCC provider such as Github.
	ErrSCC = newErr("E10010", codes.Unavailable, http.StatusServiceUnavailable, "there was an error interacting with the source code provider")
	// Means a provided connection ID was not found in the database.
	ErrConnectionNotFound = newErr("E10011", codes.NotFound, http.StatusNotFound, "connection not found")
	// Returned if an account id is not found in the database
	ErrAccountNotFound = newErr("E10012", codes.NotFound, http.StatusNotFound, "account not found")
	// Returned if an account id is not valid
	ErrInvalidAccountID = newErr("E10013", codes.InvalidArgument, http.StatusBadRequest, "invalid account id")
	// Returned if a policy id is not found in the database
	ErrPolicyNotFound = newErr("E10014", codes.NotFound, http.StatusNotFound, "policy not found")
	// Returned when there's a problem with one of the system connections
	ErrSystemConnection = newErr("E10015", codes.Internal, http.StatusInternalServerError, "system connection problem")
	// Returned if a policy id is invalid
	ErrInvalidPolicyID = newErr("E10016", codes.InvalidArgument, http.StatusBadRequest, "invalid policy id")
	// Returned when there's a problem with a connection's secret
	ErrConnectionSecret = newErr("E10017", codes.Unavailable, http.StatusInternalServerError, "connection secret error")
	// Returned when an invite for an email already exists
	ErrInviteExists = newErr("E10018", codes.AlreadyExists, http.StatusConflict, "invite already exists")
	// Returned when an invitation has expired
	ErrInviteExpired = newErr("E10019", codes.AlreadyExists, http.StatusConflict, "invite is expired")
	// Means an existing member of a tenant was invited to join the same tenant
	ErrAlreadyMember = newErr("E10020", codes.AlreadyExists, http.StatusConflict, "already a tenant member")
	// Returned if an account tried to accept or decline the invite of another account
	ErrInviteForAnotherUser = newErr("E10021", codes.PermissionDenied, http.StatusForbidden, "invite meant for another user")
	// Returned if an SCC repository has already been referenced in a policy
	ErrRepoAlreadyConnected = newErr("E10022", codes.AlreadyExists, http.StatusConflict, "repo has already been connected to a policy")
	// Returned if there was a problem setting up a Github secret
	ErrGithubSecret = newErr("E10023", codes.Unavailable, http.StatusServiceUnavailable, "failed to setup repo secret")
	// Returned if there was a problem setting up an Auth0 user
	ErrAuth0UserSetup = newErr("E10024", codes.Unavailable, http.StatusServiceUnavailable, "failed to setup user")
	// Returned if an invalid email address was used
	ErrInvalidEmail = newErr("E10025", codes.InvalidArgument, http.StatusBadRequest, "invalid email address")
	// Returned if a string doesn't look like an auth0 ID
	ErrInvalidAuth0ID = newErr("E10026", codes.InvalidArgument, http.StatusBadRequest, "invalid auth0 ID")
	// Returned when an invitation has been accepted
	ErrInviteAlreadyAccepted = newErr("E10027", codes.AlreadyExists, http.StatusConflict, "invite has already been accepted")
	// Returned when an invitation has been declined
	ErrInviteAlreadyDeclined = newErr("E10028", codes.AlreadyExists, http.StatusConflict, "invite has already been declined")
	// Returned when an invitation has been canceled
	ErrInviteCanceled = newErr("E10029", codes.AlreadyExists, http.StatusConflict, "invite has been canceled")
	// Returned when a provider verification call has failed
	ErrProviderVerification = newErr("E10030", codes.InvalidArgument, http.StatusBadRequest, "verification failed")
	// Means an account already exists for the specified user
	ErrHasAccount = newErr("E10031", codes.AlreadyExists, http.StatusConflict, "already has an account")
	// Returned when a user is not allowed to perform an operation
	ErrNotAllowed = newErr("E10032", codes.PermissionDenied, http.StatusForbidden, "not allowed")
	// Returned when trying to delete the last owner of a tenant
	ErrLastOwner = newErr("E10033", codes.PermissionDenied, http.StatusForbidden, "last owner of the tenant")
	// Returned when an operation timed out after multiple retries
	ErrRetryTimeout = newErr("E10034", codes.DeadlineExceeded, http.StatusRequestTimeout, "timeout after multiple retries")
	// Returned when a field is marked as an ID, and it's not a string
	ErrInvalidIDType = newErr("E10035", codes.InvalidArgument, http.StatusBadRequest, "ID fields have to be strings")
	// Returned when an ID is not correct
	ErrInvalidID = newErr("E10036", codes.InvalidArgument, http.StatusBadRequest, "invalid ID type")
	// Returned when trying to delete an entity that still has dependants
	ErrNotEmpty = newErr("E10037", codes.FailedPrecondition, http.StatusBadRequest, "entity is not empty")
	// Returned when authentication has failed or is not possible
	ErrAuthenticationFailed = newErr("E10038", codes.FailedPrecondition, http.StatusUnauthorized, "authentication failed")
	// Returned when a given parameter is incorrect (wrong format, value or type)
	ErrInvalidArgument = newErr("E10039", codes.InvalidArgument, http.StatusBadRequest, "invalid argument")
	// Returned when the caller is trying to update a readonly value
	ErrReadOnly = newErr("E10040", codes.InvalidArgument, http.StatusBadRequest, "readonly")
	// Returned when the caller tries to create or update a policy with a name that already exists
	ErrDuplicatePolicyName = newErr("E10041", codes.InvalidArgument, http.StatusConflict, "policy name already exists")
	// Returned when the caller tries to create or update a connection with a name that already exists
	ErrDuplicateConnectionName = newErr("E10042", codes.InvalidArgument, http.StatusConflict, "connection name already exists")
	// Returned if a module is not found
	ErrModuleNotFound = newErr("E10043", codes.NotFound, http.StatusNotFound, "module not found")
	// Return if a user is not found
	ErrUserNotFound = newErr("E10044", codes.NotFound, http.StatusNotFound, "user not found")
	// Return if a user already exists
	ErrUserAlreadyExists = newErr("E10045", codes.AlreadyExists, http.StatusConflict, "user already exists")
	// Returned when authorization has failed or is not possible
	ErrAuthorizationFailed = newErr("E10046", codes.PermissionDenied, http.StatusUnauthorized, "authorization failed")
	// Returned when a runtime query has an error
	ErrBadQuery = newErr("E10047", codes.InvalidArgument, http.StatusBadRequest, "invalid query")
	// Returned when a runtime query has an error
	ErrQueryExecutionFailed = newErr("E10048", codes.FailedPrecondition, http.StatusBadRequest, "query failed")
	// Returned when the account has not setup a personal tenant yet
	ErrPersonalTenantRequired = newErr("E10049", codes.FailedPrecondition, http.StatusBadRequest, "personal tenant required")

	asertoErrors = make(map[string]*AsertoError)
)

func newErr(code string, statusCode codes.Code, httpCode int, msg string) *AsertoError {
	asertoError := &AsertoError{code, statusCode, msg, httpCode, map[string]string{}, nil}
	asertoErrors[code] = asertoError
	return asertoError
}

// AsertoError represents a well known error
// comming from an Aserto service
type AsertoError struct {
	Code       string
	StatusCode codes.Code
	Message    string
	HttpCode   int
	data       map[string]string
	errs       []error
}

func (e *AsertoError) Data() map[string]string {
	return e.Copy().data
}

// SameAs returns true if the provided error is an AsertoError
// and has the same error code
func (e *AsertoError) SameAs(err error) bool {
	aErr, ok := err.(*AsertoError)
	if !ok {
		return false
	}

	return aErr.Code == e.Code
}

func (e *AsertoError) Copy() *AsertoError {
	dataCopy := make(map[string]string, len(e.data))

	for k, v := range e.data {
		dataCopy[k] = v
	}

	return &AsertoError{
		Code:       e.Code,
		StatusCode: e.StatusCode,
		Message:    e.Message,
		data:       dataCopy,
		errs:       e.errs,
		HttpCode:   e.HttpCode,
	}
}

func (e *AsertoError) Error() string {
	innerMessage := ""

	if len(e.errs) > 0 {
		innerMessage = e.errs[0].Error()

		for _, err := range e.errs[1:] {
			innerMessage = strings.Join([]string{innerMessage, err.Error()}, ": ")
		}
	}

	if innerMessage == "" {
		return fmt.Sprintf("%s %s", e.Code, e.Message)
	} else {
		return fmt.Sprintf("%s %s: %s", e.Code, e.Message, innerMessage)
	}
}

func (e *AsertoError) Fields() map[string]interface{} {
	result := make(map[string]interface{}, len(e.data))

	for k, v := range e.data {
		result[k] = v
	}

	return result
}

// Associates err with the AsertoError.
func (e *AsertoError) Err(err error) *AsertoError {
	if err == nil {
		return e
	}
	c := e.Copy()

	c.errs = append(c.errs, err)

	if aErr, ok := err.(*AsertoError); ok {
		for k, v := range aErr.data {
			if _, ok := c.data[k]; !ok {
				c.data[k] = v
			}
		}
	}

	return c
}

func (e *AsertoError) Msg(message string) *AsertoError {
	c := e.Copy()

	if existingMsg, ok := c.data[MessageKey]; ok {
		c.data[MessageKey] = strings.Join([]string{existingMsg, message}, ": ")
	} else {
		c.data[MessageKey] = message
	}

	return c
}

func (e *AsertoError) Msgf(message string, args ...interface{}) *AsertoError {
	c := e.Copy()

	message = fmt.Sprintf(message, args...)

	if existingMsg, ok := c.data[MessageKey]; ok {
		c.data[MessageKey] = strings.Join([]string{existingMsg, message}, ": ")
	} else {
		c.data[MessageKey] = message
	}
	return c
}

func (e *AsertoError) Str(key, value string) *AsertoError {
	c := e.Copy()
	c.data[key] = value
	return c
}

func (e *AsertoError) Int(key string, value int) *AsertoError {
	c := e.Copy()
	c.data[key] = fmt.Sprintf("%d", value)
	return c
}

func (e *AsertoError) Bool(key string, value bool) *AsertoError {
	c := e.Copy()
	c.data[key] = fmt.Sprintf("%t", value)

	return c
}

func (e *AsertoError) Duration(key string, value time.Duration) *AsertoError {
	c := e.Copy()
	c.data[key] = value.String()
	return c
}

func (e *AsertoError) Time(key string, value time.Time) *AsertoError {
	c := e.Copy()
	c.data[key] = value.UTC().Format(time.RFC3339)
	return c
}

func (e *AsertoError) FromReader(key string, value io.Reader) *AsertoError {
	buf := &strings.Builder{}
	_, err := io.Copy(buf, value)

	if err != nil {
		return e.Err(err)
	}

	c := e.Copy()
	c.data[key] = buf.String()

	return c
}

func (e *AsertoError) Interface(key string, value interface{}) *AsertoError {
	c := e.Copy()
	c.data[key] = fmt.Sprintf("%+v", value)
	return c
}

func (e *AsertoError) Unwrap() error {
	if len(e.errs) > 0 {
		return e.errs[len(e.errs)-1]
	}

	return nil
}

func (e *AsertoError) Cause() error {
	if len(e.errs) > 0 {
		return e.errs[len(e.errs)-1]
	}

	return nil
}

func (e *AsertoError) MarshalZerologObject(event *zerolog.Event) {
	event.Str("error", e.Error())
	event.Fields(e.Fields())
}

func (e *AsertoError) GRPCStatus() *status.Status {
	return status.New(e.StatusCode, e.Message)
}

func (e *AsertoError) WithGRPCStatus(grpcCode codes.Code) *AsertoError {
	c := e.Copy()
	c.StatusCode = grpcCode
	return c
}

func (e *AsertoError) WithHTTPStatus(httpStatus int) *AsertoError {
	c := e.Copy()
	c.HttpCode = httpStatus
	return c
}

// Returns an Aserto error based on a given grpcStatus. The details that are not of type errdetails.ErrorInfo are dropped.
// and if there are details from multiple errors, the aserto error will be constructed based on the first one.
func FromGRPCStatus(grpcStatus status.Status) *AsertoError {
	var result *AsertoError
	for _, detail := range grpcStatus.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			result = asertoErrors[t.Domain]
			result.data = t.Metadata
		}
		if result != nil {
			break
		}
	}
	return result
}

func UnwrapAsertoError(err error) *AsertoError {
	initialError := errors.Cause(err)
	if initialError == nil {
		initialError = err
	}
	grpcStatus, ok := status.FromError(initialError)
	if ok {
		aErr := FromGRPCStatus(*grpcStatus)
		if aErr != nil {
			return aErr
		}
	}

	for {
		aErr, ok := err.(*AsertoError)
		if ok {
			return aErr
		}

		err = errors.Unwrap(err)
		if err == nil {
			return nil
		}
	}
}
