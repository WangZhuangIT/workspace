
package errorutil

type ERROR int

const (
	ERROR_NONE                = 200
	ERROR_NOT_MODIFIED        = 304 // RFC 7232, 4.1
	ERROR_CLIENT              = 400
	ERROR_NEED_LOGIN          = 401
	ERROR_NOT_PERMIT          = 403
	ERROR_NOT_FOUND           = 404
	ERROR_SERVER              = 500
	ERROR_PANIC               = 502
	ERROR_INVALID_PARAM ERROR = 401000 + iota
)

