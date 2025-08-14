package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

	testCases := []struct {
		name          string
		accountid     int64
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountid: account.ID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatching(t, recorder.Body, account)

			},
		},

		{
			name:      "NotFound",
			accountid: account.ID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// removed RequireBodyMatching() because there is no body to match or test
			},
		},


		{
			name:      "BadRequest",
			accountid: 0,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				

			},
		},

		{
			name:      "InternalServerError",
			accountid: account.ID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			//It starts test Server and sends Requests(like GetAccount)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomAccount() Anuskh.Account {
	return Anuskh.Account{
		ID:       util.RandomInt(1, 100),
		Balance:  util.RandomBalance(),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}
}

func RequireBodyMatching(t *testing.T, body *bytes.Buffer, account Anuskh.Account) {

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotAccount Anuskh.Account

	err = json.Unmarshal(data, &GotAccount)
	require.NoError(t, err)

	require.Equal(t, account, GotAccount)
}


