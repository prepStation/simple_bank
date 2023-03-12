package token

import (
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {

	//CreteToken creates a new token for a specific username and duration
	CreateToken(userName string, duration time.Duration) (string, error)

	//VerifyToken verifies if the input token is valid or not
	VerifyToken(token string) (*Payload, error)
}
