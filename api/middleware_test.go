package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	token "simplebank/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(t *testing.T, request *http.Request, tokenMaker token.Maker, authType string, email string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set("Authorization", authorizationHeader)

}
func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, "email", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NOAUTH",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		}, {
			name: "UnsuportedAUTH",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, "unsupported", "email", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthFormat",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, "", "email", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		}, {
			name: "UnsuportedAUTH",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, "email", -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tcase := testCases[i]
		t.Run(tcase.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, AuthMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)

			require.NoError(t, err)

			tcase.setUpAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tcase.checkResponse(t, recorder)

		})

	}
}
