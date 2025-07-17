package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/model"
	"simplebank/factory"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateOwner(t *testing.T) {
	owner := randomOwner()

	/*reqBody := map[string]interface{}{
		"first_name":     owner.FirstName,
		"first_surname":  owner.FirstSurname,
		"second_surname": owner.SecondSurname,
		"born_at":        owner.BornAt.Format(time.DateOnly),
		"nationality":    owner.Nationality,
		"email":          owner.Email,
		"password":       "testPassword",
	}*/

	log.Println(owner.BornAt.Format("2006-01-02"))
	log.Println(owner.BornAt.Format("2006-01-02"))

	testCases := []struct {
		name          string
		createBody    func(owner *db.Owner) map[string]interface{}
		buildStbus    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOwner(gomock.Any(), gomock.AssignableToTypeOf(db.CreateOwnerParams{})).
					Times(1).
					Return(owner, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)

			},
			createBody: func(owner *db.Owner) map[string]interface{} {
				return map[string]interface{}{
					"first_name":     owner.FirstName,
					"first_surname":  owner.FirstSurname,
					"second_surname": owner.SecondSurname,
					"born_at":        owner.BornAt.Format(time.DateOnly),
					"nationality":    owner.Nationality,
					"email":          owner.Email,
					"password":       "testPassword",
				}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOwner(t, recorder.Body, owner)
			},
		},
		{
			name: "Bad Request",
			createBody: func(owner *db.Owner) map[string]interface{} {
				return map[string]interface{}{
					"first_namEEEEEEEEEEREREWREREW": owner.FirstName,
					"first_surname":                 owner.FirstSurname,
					"second_surname":                owner.SecondSurname,
					"born_at":                       owner.BornAt.Format(time.DateOnly),
					"nationality":                   owner.Nationality,
					"email":                         owner.Email,
					"password":                      "testPassword",
				}
			},
			buildStbus: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOwner(gomock.Any(), gomock.AssignableToTypeOf(db.CreateOwnerParams{})).
					Times(0) //.
				//Return(db.Owner{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//requireBodyMatchOwner(t, recorder.Body, owner)
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

			url := "/owner"

			reqBody := tc.createBody(&owner)

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
}

func requireBodyMatchOwner(t *testing.T, body *bytes.Buffer, owner db.Owner) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var rqOwner db.Owner
	err = json.Unmarshal(data, &rqOwner)
	require.NoError(t, err)

	require.Equal(t, rqOwner.FirstName, owner.FirstName)
	require.Equal(t, rqOwner.FirstSurname, owner.FirstSurname)
	require.Equal(t, rqOwner.SecondSurname, owner.SecondSurname)
	require.Equal(t, rqOwner.Email, owner.Email)
	require.Equal(t, rqOwner.Nationality, owner.Nationality)

	// Compara només la data (no l’hora)
	require.Equal(t, rqOwner.BornAt.Format("2006-01-02"), owner.BornAt.Format("2006-01-02"))

}

func randomOwner() db.Owner {

	return db.Owner{
		ID:             factory.RandomInt(1, 1000),
		FirstName:      factory.RandomString(7),
		FirstSurname:   factory.RandomString(7),
		SecondSurname:  factory.RandomString(8),
		BornAt:         time.Now(),
		Nationality:    int32(factory.RandomInt(1, 99)),
		HashedPassword: factory.RandomString(8),
		Email:          factory.RandomEmail(),
	}
}
