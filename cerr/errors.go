package cerr

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MessageKey = "msg"
)

var (
	// Unknown error ID. It's returned when the implementation has not returned another AsertoError.
	ErrUnknown = newErr("E10000", codes.Internal, "an unknown error has occurred")
	// Means no tenant id was found in the current context
	ErrNoTenantID = newErr("E10001", codes.InvalidArgument, "no tenant id specified")
	// Means the tenant id is not valid
	ErrInvalidTenantID = newErr("E10002", codes.InvalidArgument, "invalid tenant id")
	// Means the tenant name doesn't conform to our tenant name rules
	ErrInvalidTenantName = newErr("E10003", codes.InvalidArgument, "invalid tenant name")
	// Means the provider ID is invalid
	ErrInvalidProviderID = newErr("E10004", codes.InvalidArgument, "invalid provider id")
	// Means the provider config name doesn't exist
	ErrInvalidProviderConfigName = newErr("E10005", codes.InvalidArgument, "invalid provider config name")
	// The asked-for runtime is not yet available, but will likely be in the future.
	ErrRuntimeLoading = newErr("E10006", codes.Unavailable, "runtime has not yet loaded")
	// Means a connection failed to validate.
	ErrConnectionVerification = newErr("E10007", codes.FailedPrecondition, "connection verification failed")
	// Returned when there's a problem retrieving a connection.
	ErrConnection = newErr("E10008", codes.Unavailable, "connection problem")
	// Returned when there's a problem getting a github access token.
	ErrGithubAccessToken = newErr("E10009", codes.Unavailable, "failed to retrieve github access token")
	// Returned when there's a problem communicating with an SCC provider such as Github.
	ErrSCC = newErr("E10010", codes.Unavailable, "there was an error interacting with the source code provider")
	// Means a provided connection ID was not found in the database.
	ErrConnectionNotFound = newErr("E10011", codes.NotFound, "connection not found")
	// Returned if an account id is not found in the database
	ErrAccountNotFound = newErr("E10012", codes.NotFound, "account not found")
	// Returned if an account id is not valid
	ErrInvalidAccountID = newErr("E10013", codes.InvalidArgument, "invalid account id")
	// Returned if a policy id is not found in the database
	ErrPolicyNotFound = newErr("E10014", codes.NotFound, "policy not found")
	// Returned when there's a problem with one of the system connections
	ErrSystemConnection = newErr("E10015", codes.Internal, "system connection problem")
	// Returned if a policy id is invalid
	ErrInvalidPolicyID = newErr("E10016", codes.InvalidArgument, "invalid policy id")
	// Returned when there's a problem with a connection's secret
	ErrConnectionSecret = newErr("E10017", codes.Unavailable, "connection secret error")
	// Returned when an invite for an email already exists
	ErrInviteExists = newErr("E10018", codes.AlreadyExists, "invite already exists")
	// Returned when an invitation has expired
	ErrInviteExpired = newErr("E10019", codes.AlreadyExists, "invite is expired")
	// Means an existing member of a tenant was invited to join the same tenant
	ErrAlreadyMember = newErr("E10020", codes.AlreadyExists, "already a tenant member")
	// Returned if an account tried to accept or decline the invite of another account
	ErrInviteForAnotherUser = newErr("E10021", codes.PermissionDenied, "invite meant for another user")
	// Returned if an SCC repository has already been referenced in a policy
	ErrRepoAlreadyConnected = newErr("E10022", codes.AlreadyExists, "repo has already been connected to a policy")
	// Returned if there was a problem setting up a Github secret
	ErrGithubSecret = newErr("E10023", codes.Unavailable, "failed to setup repo secret")
	// Returned if there was a problem setting up an Auth0 user
	ErrAuth0UserSetup = newErr("E10024", codes.Unavailable, "failed to setup user")
	// Returned if an invalid email address was used
	ErrInvalidEmail = newErr("E10025", codes.InvalidArgument, "invalid email address")
	// Returned if a string doesn't look like an auth0 ID
	ErrInvalidAuth0ID = newErr("E10026", codes.InvalidArgument, "invalid auth0 ID")
	// Returned when an invitation has been accepted
	ErrInviteAlreadyAccepted = newErr("E10027", codes.AlreadyExists, "invite has already been accepted")
	// Returned when an invitation has been declined
	ErrInviteAlreadyDeclined = newErr("E10028", codes.AlreadyExists, "invite has already been declined")
	// Returned when an invitation has been canceled
	ErrInviteCanceled = newErr("E10029", codes.AlreadyExists, "invite has been canceled")
	// Returned when a provider verification call has failed
	ErrProviderVerification = newErr("E10030", codes.InvalidArgument, "verification failed")
	// Means an account already exists for the specified user
	ErrHasAccount = newErr("E10031", codes.AlreadyExists, "already has an account")
	// Returned when a user is not allowed to perform an operation
	ErrNotAllowed = newErr("E10032", codes.PermissionDenied, "not allowed")
	// Returned when trying to delete the last owner of a tenant
	ErrLastOwner = newErr("E10033", codes.PermissionDenied, "last owner of the tenant")
	// Returned when an operation timed out after multiple retries
	ErrRetryTimeout = newErr("E10034", codes.DeadlineExceeded, "timeout after multiple retries")
	// Returned when a field is marked as an ID, and it's not a string
	ErrInvalidIDType = newErr("E10035", codes.InvalidArgument, "ID fields have to be strings")
	// Returned when an ID is not correct
	ErrInvalidID = newErr("E10036", codes.InvalidArgument, "invalid ID type")
	// Returned when trying to delete an entity that still has dependants
	ErrNotEmpty = newErr("E10037", codes.FailedPrecondition, "entity is not empty")
	// Returned when authentication has failed or is not possible
	ErrAuthenticationFailed = newErr("E10038", codes.FailedPrecondition, "authentication failed")
	// Returned when a given parameter is incorrect (wrong format, value or type)
	ErrInvalidArgument = newErr("E10039", codes.InvalidArgument, "invalid argument")
	// Returned when the caller is trying to update a readonly value
	ErrReadOnly = newErr("E10040", codes.InvalidArgument, "readonly")
)

func newErr(code string, statusCode codes.Code, msg string) *AsertoError {
	return &AsertoError{code, statusCode, msg, map[string]string{}, nil}
}

// AsertoError represents a well known error
// comming from an Aserto service
type AsertoError struct {
	Code       string
	StatusCode codes.Code
	Message    string
	Data       map[string]string
	errs       []error
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
	dataCopy := make(map[string]string, len(e.Data))

	for k, v := range e.Data {
		dataCopy[k] = v
	}

	return &AsertoError{
		Code:       e.Code,
		StatusCode: e.StatusCode,
		Message:    e.Message,
		Data:       dataCopy,
		errs:       e.errs,
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
	result := make(map[string]interface{}, len(e.Data))

	for k, v := range e.Data {
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
		for k, v := range aErr.Data {
			if _, ok := c.Data[k]; !ok {
				c.Data[k] = v
			}
		}
	}

	return c
}

func (e *AsertoError) Msg(message string) *AsertoError {
	c := e.Copy()

	if existingMsg, ok := c.Data[MessageKey]; ok {
		c.Data[MessageKey] = strings.Join([]string{existingMsg, message}, ": ")
	} else {
		c.Data[MessageKey] = message
	}

	return c
}

func (e *AsertoError) Msgf(message string, args ...interface{}) *AsertoError {
	c := e.Copy()

	message = fmt.Sprintf(message, args...)

	if existingMsg, ok := c.Data[MessageKey]; ok {
		c.Data[MessageKey] = strings.Join([]string{existingMsg, message}, ": ")
	} else {
		c.Data[MessageKey] = message
	}
	return c
}

func (e *AsertoError) Str(key, value string) *AsertoError {
	c := e.Copy()
	c.Data[key] = value
	return c
}

func (e *AsertoError) Int(key string, value int) *AsertoError {
	c := e.Copy()
	c.Data[key] = fmt.Sprintf("%d", value)
	return c
}

func (e *AsertoError) Bool(key string, value bool) *AsertoError {
	c := e.Copy()
	c.Data[key] = fmt.Sprintf("%t", value)
	return c
}

func (e *AsertoError) Duration(key string, value time.Duration) *AsertoError {
	c := e.Copy()
	c.Data[key] = value.String()
	return c
}

func (e *AsertoError) Time(key string, value time.Time) *AsertoError {
	c := e.Copy()
	c.Data[key] = value.UTC().Format(time.RFC3339)
	return c
}

func (e *AsertoError) FromReader(key string, value io.Reader) *AsertoError {
	buf := &strings.Builder{}
	_, err := io.Copy(buf, value)

	if err != nil {
		return e.Err(err)
	}

	c := e.Copy()
	c.Data[key] = buf.String()

	return c
}

func (e *AsertoError) Interface(key string, value interface{}) *AsertoError {
	c := e.Copy()
	c.Data[key] = fmt.Sprintf("%+v", value)
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
