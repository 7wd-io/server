package token

import (
	"7wd.io/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken(t *testing.T) {
	f := New("secret")

	token, err := f.Token(&domain.Passport{
		Id:       1,
		Nickname: "user1",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
