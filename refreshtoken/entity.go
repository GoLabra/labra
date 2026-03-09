package refreshtoken

import (
	"database/sql"
	"errors"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type StoredRefreshToken struct {
	SubjectEmail  string
	SubjectType   string
	RoleName      string
	ExpiresAtUnix int64
	RevokedAtUnix sql.NullInt64
}
