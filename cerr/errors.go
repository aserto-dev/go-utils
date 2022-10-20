package cerr

import (
	"net/http"

	"github.com/aserto-dev/errors"
	"google.golang.org/grpc/codes"
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
	// Returned when trying to delete an entity that still has dependents
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
	// Returned if a policy builder id is not found in the database
	ErrPolicyBuilderNotFound = newErr("E10050", codes.NotFound, http.StatusNotFound, "policy builder not found")
	// Returned if a policy builder id is invalid
	ErrInvalidPolicyBuilderID = newErr("E10051", codes.InvalidArgument, http.StatusBadRequest, "invalid policy builder id")
	// Returned when discovery for policy runtime configuration has failed
	ErrDiscoveryFailed = newErr("E10051", codes.Unavailable, http.StatusServiceUnavailable, "discovery failed")
	// Returned when a decision is invalid
	ErrInvalidDecision = newErr("E10052", codes.InvalidArgument, http.StatusBadRequest, "invalid decision")
	// Returned when a runtime failed to load
	ErrBadRuntime = newErr("E10053", codes.Unavailable, http.StatusServiceUnavailable, "runtime loading failed")
	// Returned when there's a problem getting a gitlab access token.
	ErrGitlabAccessToken = newErr("E10054", codes.Unavailable, http.StatusServiceUnavailable, "failed to retrieve gitlab access token")
	// Returned when a tag is not a valid policy
	ErrInvalidPolicyTag = newErr("E10055", codes.InvalidArgument, http.StatusBadRequest, "invalid policy tag")
	// Returned if a policy instance is not found in the database
	ErrPolicyInstanceNotFound = newErr("E10056", codes.NotFound, http.StatusNotFound, "policy instance not found")
	// Returned if a policy repository is not found in the database
	ErrPolicyRepositoryNotFound = newErr("E10057", codes.NotFound, http.StatusNotFound, "policy repository not found")
	// Returned if a policy source is not found in the database
	ErrPolicySourceNotFound = newErr("E10058", codes.NotFound, http.StatusNotFound, "policy source not found")
	// Returned if a source has already been attached to a policy
	ErrPolicySourceAlreadySet = newErr("E10059", codes.FailedPrecondition, http.StatusBadRequest, "source already set")
	// Returned if the organization is not found in the source code provider
	ErrSCCOrganizationNotFound = newErr("E10060", codes.NotFound, http.StatusNotFound, "source code control organization not found")
	// Returned if the repo is not found in the source code provider
	ErrSCCRepoNotFound = newErr("E10061", codes.NotFound, http.StatusNotFound, "source code control repository not found")
	// Returned if a policy already has a connected repository
	ErrPolicyRepositoryAlreadyConnected = newErr("E10062", codes.AlreadyExists, http.StatusConflict, "the policy already has a repository connected")
	// Returned if object type is not defined in the directory
	ErrDirectoryObjectTypeUnknown = newErr("E10063", codes.Unknown, http.StatusNotFound, "directory object type unknown")
	// Returned if relation type is not defined in the directory
	ErrDirectoryRelationTypeUnknown = newErr("E10064", codes.NotFound, http.StatusNotFound, "directory relation type unknown")
	// Returned if permission is not defined in the directory
	ErrDirectoryPermissionUnknown = newErr("E10065", codes.NotFound, http.StatusNotFound, "directory permission unknown")
	// Returned if object object id is not found in the directory
	ErrDirectoryObjectNotFound = newErr("E10066", codes.NotFound, http.StatusNotFound, "directory object not found")
	// Returned if relation object is not found in the directory
	ErrDirectoryRelationNotFound = newErr("E10067", codes.NotFound, http.StatusNotFound, "directory relation not found")
	// Returned if the tenant is marked for deletion
	ErrTenantDeleted = newErr("E10068", codes.NotFound, http.StatusNotFound, "tenant is marked for deletion")
	// Returned when tenant store for given tenant id is not found in directory.
	ErrDirectoryStoreTenantNotFound = newErr("E10069", codes.NotFound, http.StatusNotFound, "tenant store not found")
	// Returned when trying to update a resource that was changed in the meanwhile
	ErrVersionsMismatch = newErr("E10070", codes.FailedPrecondition, http.StatusPreconditionFailed, "version hash mismatch")
	// Returned if a tenant id is not found in the database
	ErrTenantNotFound = newErr("E10071", codes.NotFound, http.StatusNotFound, "tenant not found")
)

func newErr(code string, statusCode codes.Code, httpCode int, msg string) *errors.AsertoError {
	return errors.NewAsertoError(code, statusCode, httpCode, msg)
}
