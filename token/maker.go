package token

import "time"

type Maker interface {
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)
	Verify(token string) (*Payload, error)
}
