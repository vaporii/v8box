package security

import "github.com/google/uuid"

func GenerateStateToken() string {
	return uuid.NewString()
}
