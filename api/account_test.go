package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/model"
	"simplebank/factory"
	token "simplebank/token"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	owner := randomOwner()
	account := randomAcount(owner.ID)
	unauthorizedOwner := randomOwner()

	testCases := []struct {
		name      string
		accountID int64
		setUpAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)

		buildStbus    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, owner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)

				store.EXPECT().
					GetOwnerByEmail(gomock.Any(), gomock.Eq(owner.Email)).
					Times(1).
					Return(owner, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NOTFOUND",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, owner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				//No cal comparar ja que no ens el troba
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, owner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, owner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Unauthized",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, unauthorizedOwner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)

				store.EXPECT().
					GetOwnerByEmail(gomock.Any(), gomock.Eq(unauthorizedOwner.Email)).
					Times(1).
					Return(unauthorizedOwner, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		{
			name:      "NoAuth",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//addAuth(t, request, tokenMaker, authTypeBearer, owner.Email, time.Minute)
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0).
					Return(account, nil)

				store.EXPECT().
					GetOwnerByEmail(gomock.Any(), gomock.Eq(owner.Email)).
					Times(0).
					Return(owner, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStbus(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setUpAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
	/*ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// els stubs
	/*store.EXPECT().
	GetAccount(gomock.Any(), gomock.Eq(account.ID)).
	Times(1).
	Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// aqui checkegem el status
	//require.Equal(t, http.StatusOK, recorder.Code)

	//requireBodyMatchAccount(t, recorder.Body, account)*/

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rqAccount db.Account
	err = json.Unmarshal(data, &rqAccount)
	require.NoError(t, err)

	require.Equal(t, rqAccount, account)

}

func randomAcount(ownerId int64) db.Account {
	return db.Account{
		ID:          factory.RandomInt(1, 1000),
		OwnerID:     ownerId, //factory.RandomOwner(),
		Currency:    "USD",
		Money:       factory.RandomMoney(),
		CountryCode: int32(factory.RandomInt(1, 32)),
	}
}
