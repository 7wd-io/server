package password

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	m := New()

	hashed, err := m.Hash("12345678", 10)

	fmt.Println(hashed)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
}
