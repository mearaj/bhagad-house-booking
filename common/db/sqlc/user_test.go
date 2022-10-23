package sqlc

import (
	"context"
	"database/sql"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func createRandomUser(t *testing.T) User {
	name := utils.RandomName()
	pass := utils.RandomString(6)
	log.Println(name)
	log.Println(pass)
	hashedPassword, err := utils.HashPassword(pass)
	require.NoError(t, err)
	arg := CreateUserParams{
		Name:     name,
		Email:    utils.RandomEmail(),
		Password: hashedPassword,
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Name, arg.Name)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.Password, arg.Password)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.Email, user2.Email)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	arg := UpdateUserParams{
		ID:    user.ID,
		Name:  "Rahim",
		Email: utils.RandomEmail(),
	}
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, arg.Name, user2.Name)
	require.Equal(t, arg.Email, user2.Email)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
	user2, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
