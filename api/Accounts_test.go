package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockDB "github.com/nilesh0729/OrdinaryBank/db/Mock"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
	"github.com/nilesh0729/OrdinaryBank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockDB.NewMockStore(ctrl)

	store.EXPECT().
		GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	//It starts test Server and sends Requests(like GetAccount)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response

	require.Equal(t, http.StatusOK, recorder.Code)

}

func randomAccount() Anuskh.Account {
	return Anuskh.Account{
		ID:       util.RandomInt(1, 100),
		Balance:  util.RandomBalance(),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}
}
