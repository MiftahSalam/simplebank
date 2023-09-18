package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simplebank/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(
	t *testing.T,
	request *http.Request,
	tokenManager token.TokenManager,
	authType string,
	username string,
	duration time.Duration) {
	token, _, err := tokenManager.CreateToken(username, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(authHeaderKey, authHeader)
}

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenManager token.TokenManager)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuth(t, request, tokenManager, authHeaderTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Authorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported Auth Type",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuth(t, request, tokenManager, "authHeaderTypeBearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Auth format",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuth(t, request, tokenManager, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager) {
				addAuth(t, request, tokenManager, authHeaderTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authpath := "/auth"
			server := newTestServer(t, nil)
			server.router.GET(authpath, authMiddleware(server.tokenManager), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authpath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
