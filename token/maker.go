package token

import "time"

// It's a Token Maker interface was implemented by JWT and PASETO token
// Easily switch between token maker.
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not
	// returns the Payload data store in the body of the token.
	VerifyToken(token string) (*Payload, error)
}
