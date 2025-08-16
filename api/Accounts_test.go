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

	"github.com/gin-gonic/gin"
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

func TestPostAccountAPI(t *testing.T) {

	account1 := randomAccount()

	arg := Anuskh.CreateAccountsParams{
		Owner:    account1.Owner,
		Currency: account1.Currency,
		Balance:  0,
	}

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockDB.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"owner":    account1.Owner,
				"Currency": account1.Currency,
				"Balance":  account1.Balance,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account1, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatching(t, recorder.Body, account1)
			},
		},

		{
			name: "BadRequest",
			body: gin.H{
				"owner":    util.RandomOwner(),
				"Currency": util.RandomBalance(),
				"Balance":  util.RandomInt(1, 100),
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name: "InternalServerError",
			body: gin.H{
				"owner":    account1.Owner,
				"Currency": account1.Currency,
				"Balance":  0,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testcases {

		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			Server := NewServer(store)

			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			url := "/accounts"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			body := bytes.NewBuffer(data)
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)

			Server.router.ServeHTTP(recorder, request)

			tc.CheckResponse(t, recorder)
		})
	}

}

func TestListAccountAPI(t *testing.T) {
	accountList := []Anuskh.Account{
		randomAccount(),
		randomAccount(),
	}
	arg := Anuskh.ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	testcases := []struct {
		name          string
		query         func(pageID, PageSize int) string
		buildStubs    func(store *mockDB.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			query: func(pageID, PageSize int) string {
				return fmt.Sprintf("page_id=%d&page_size=%d", pageID, PageSize)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accountList, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatchingAccounts(t, recorder.Body, accountList)
			},
		},

		{
			name: "BadRequest",
			query: func(pageID, PageSize int) string {
				return fmt.Sprintf("page_id=%d&page_size=%d", 0, PageSize)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
					
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				
			},
		},

		{
			name: "InternalServerError",
			query: func(pageID, PageSize int) string {
				return fmt.Sprintf("page_id=%d&page_size=%d", pageID, PageSize)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]Anuskh.Account{}, sql.ErrConnDone)
					
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				
			},
		},
		
	}

	for i := range testcases {
		
		tc := testcases[i]
		t.Run(tc.name,func(t *testing.T){

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockDB.NewMockStore(ctrl)
		tc.buildStubs(store)

		Server := NewServer(store)

		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts?%s", tc.query(1, 5))

		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		Server.router.ServeHTTP(recorder, request)

		tc.CheckResponse(t, recorder)
		})
	}

}

func RequireBodyMatchingAccounts(t *testing.T, body *bytes.Buffer, expected []Anuskh.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got []Anuskh.Account
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Equal(t, expected, got)
}
