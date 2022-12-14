// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: account.sql

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jiny0x01/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users{
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}	

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}


func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
				String: newFullName,
				Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, newFullName)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
				String: newEmail,
				Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, newEmail)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.FullName, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newHashPassword, _ := util.HashPassword(util.RandomString(6))
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
				String: newHashPassword,
				Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, newHashPassword)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
}

func TestUpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newHashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid: true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
		HashedPassword: sql.NullString{
			String: newHashPassword,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.NotEqual(t, oldUser.Email, newEmail)
	require.NotEqual(t, oldUser.FullName, newFullName)
	require.NotEqual(t, oldUser.HashedPassword, newHashPassword)
}