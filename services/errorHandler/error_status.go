package errorhandler

// CustomStatus type with int as the underlying type
type CustomStatus int

// Constants representing different HTTP status codes
const (
	BadRequest CustomStatus = iota + 1
	Unauthorized
	PaymentRequired
	Forbidden
	NotFound
	MethodNotAllowed
	NotAcceptable
	ProxyAuthRequired
	RequestTimeout
	Conflict
	Gone
	LengthRequired
	PreconditionFailed
	RequestEntityTooLarge
	RequestURITooLong
	UnsupportedMediaType
	RequestedRangeNotSatisfiable
	ExpectationFailed
	Teapot
	MisdirectedRequest
	UnprocessableEntity
	Locked
	FailedDependency
	UpgradeRequired
	PreconditionRequired
	TooManyRequests
	RequestHeaderFieldsTooLarge
	UnavailableForLegalReasons
	InternalServerError
	NotImplemented
	BadGateway
	ServiceUnavailable
	GatewayTimeout
	HTTPVersionNotSupported
)
