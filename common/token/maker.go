package token

import (
	"time"
)

type Maker interface {
	CreateToken(uniqueStr string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
