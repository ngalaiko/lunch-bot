package keys

import (
	"fmt"

	"github.com/google/uuid"
)

// Known errors
var (
	ErrKeyIsEmpty = fmt.Errorf("key is empty")
)

// Key is a public key.
type Key struct {
	ID        string `dynamodbav:"id" json:"id"`
	PublicDER []byte `dynamodbav:"public_der" json:"public_der"`
}

// New creates a new key with a public der payload.
func New(publicDER []byte) (*Key, error) {
	if len(publicDER) == 0 {
		return nil, ErrKeyIsEmpty
	}

	return &Key{
		ID:        uuid.New().String(),
		PublicDER: publicDER,
	}, nil
}
