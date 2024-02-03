package responsehandler

// SuccessStatus type with int as the underlying type
type SuccessStatus int

// Constants representing different HTTP status codes
const (
	OK SuccessStatus = iota
	Created
	Accepted
	NonAuthoritativeInfo
	NoContent
	ResetContent
	PartialContent
	MultiStatus
	AlreadyReported
	IMUsed
)
