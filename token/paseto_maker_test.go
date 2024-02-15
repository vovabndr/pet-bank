package token

import (
	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
	"pet-bank/utils"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := utils.RandomString(chacha20poly1305.KeySize)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := utils.RandomOwner()

	token, payload, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.Verify(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
}
