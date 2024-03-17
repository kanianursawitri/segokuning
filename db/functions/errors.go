package functions

import "errors"

var (
	ErrNoRow          = errors.New("NOT_FOUND")
	ErrInsuficientQty = errors.New("INSUFICCIENT_QUANTITY")
	ErrUnauthorized   = errors.New("UNAUTHORIZED")
	ErrDuplicate      = errors.New("DATA_DUPLICATE")
	ErrDataExists     = errors.New("DATA_EXISTS")
)
