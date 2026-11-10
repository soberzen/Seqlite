package errs

type Kind string

const (
	// 500 series errors
	KindInternal Kind = "internal" // 500
	// 400 series errors
	KindInvalid      Kind = "invalid"      // 400 bad request / 参数错误
	KindUnauthorized Kind = "unauthorized" // 401
	KindForbidden    Kind = "forbidden"    // 403
	KindNotFound     Kind = "not_found"    // 404
	KindConflict     Kind = "conflict"     // 409
)

func (k Kind) HTTPStatus() int {
	switch k {
	case KindInvalid:
		return 400
	case KindUnauthorized:
		return 401
	case KindForbidden:
		return 403
	case KindNotFound:
		return 404
	case KindConflict:
		return 409
	default:
		return 500
	}
}
