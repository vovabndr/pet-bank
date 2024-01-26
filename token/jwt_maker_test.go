package token

import (
	"github.com/stretchr/testify/require"
	"pet-bank/utils"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	secretKey := utils.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	secretKey := utils.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	token, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
