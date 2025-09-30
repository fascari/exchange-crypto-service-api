package createaccount_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"exchange-crypto-service-api/internal/app/account/handler/createaccount"
	"exchange-crypto-service-api/internal/app/account/handler/createaccount/testdata"
	createaccountuc "exchange-crypto-service-api/internal/app/account/usecase/createaccount"
	"exchange-crypto-service-api/internal/app/account/usecase/createaccount/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Handle_Success(t *testing.T) {
	mockAccountRepo := mocks.NewRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockExchangeRepo := mocks.NewExchangeRepository(t)

	mockUserRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)
	mockAccountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(testdata.ValidAccount(), nil)
	mockExchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)

	useCase := createaccountuc.New(mockAccountRepo, mockUserRepo, mockExchangeRepo)
	handler := createaccount.NewHandler(useCase)

	req := createRequest(testdata.ValidInputPayload())
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	assertSuccessResponse(t, rr, http.StatusCreated)
}

func TestHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository)
		requestBody    any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "should return 400 when invalid JSON is provided",
			setupMocks:     func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository) {},
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid payload",
		},
		{
			name:           "should return 400 when required fields are missing",
			setupMocks:     func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository) {},
			requestBody:    testdata.InputPayloadMissingUserID(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return 400 when user has invalid minimum age",
			setupMocks: func(_ *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.UnderageUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)
			},
			requestBody:    testdata.ValidInputPayload(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "user does not meet minimum age requirement",
		},
		{
			name: "should return 500 when internal error occurs",
			setupMocks: func(accountRepo *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)
				accountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(testdata.ZeroBalanceAccount(), errors.New("database connection failed"))
			},
			requestBody:    testdata.ValidInputPayload(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupHandler(t, tt.setupMocks)
			req := createRequest(tt.requestBody)
			rr := httptest.NewRecorder()

			handler.Handle(rr, req)

			assertErrorResponse(t, rr, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func setupHandler(t *testing.T, setupMocks func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository)) createaccount.Handler {
	mockAccountRepo := mocks.NewRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockExchangeRepo := mocks.NewExchangeRepository(t)

	setupMocks(mockAccountRepo, mockUserRepo, mockExchangeRepo)

	useCase := createaccountuc.New(mockAccountRepo, mockUserRepo, mockExchangeRepo)
	return createaccount.NewHandler(useCase)
}

func createRequest(requestBody any) *http.Request {
	var body []byte
	if str, ok := requestBody.(string); ok {
		body = []byte(str)
	}

	if body == nil {
		body, _ = json.Marshal(requestBody)
	}

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertSuccessResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(t, expectedStatus, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response createaccount.OutputPayload
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

	expectedOutput := createaccount.ToOutputPayload(testdata.ValidAccount())
	require.Equal(t, expectedOutput, response)
}

func assertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	require.Equal(t, expectedStatus, rr.Code)

	if expectedBody != "" {
		require.Contains(t, rr.Body.String(), expectedBody)
	}
}
