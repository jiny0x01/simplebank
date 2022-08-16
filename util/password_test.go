package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// HashPassword에서 salt할 때 랜덤으로 salt하기 때문에 같은 패스워드로 입력해도 해시 충돌이 일어나지 않음.
	// hashedPassword와 사용자가 입력한 password를 비교할 때 hashedPassword의 cost값을 구해서 그 값을 바탕으로 사용자 입력 password를 해싱하여 비교함.
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword2, hashedPassword1)
}


