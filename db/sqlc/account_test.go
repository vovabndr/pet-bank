package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"pet-bank/utils"
	"testing"
)

func createTestAccount(t *testing.T) (account Account) {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createTestAccount(t)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1, account2)
}

func TestUpdateAccountBalance(t *testing.T) {
	account1 := createTestAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account1Updated, err := testStore.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account1Updated)

	require.Equal(t, account1.ID, account1Updated.ID)
	require.Equal(t, arg.Balance, account1Updated.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createTestAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
	require.Equal(t, err, ErrRecordNotFound)
}

func TestListAccounts(t *testing.T) {
	account := createTestAccount(t)
	arg := ListAccountsParams{
		Owner:  account.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testStore.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}
