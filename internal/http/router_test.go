package http_test

/*
func TestNewRouter(t *testing.T) {
	for _, tt := range []struct {
		path               string
		method             string
		body               string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			path:               "/api/v1/exchange",
			method:             http.MethodPost,
			body:               `{"primaryCurrency":"ABC","secondaryCurrency":"XYZ","amount":100000.99}`,
			expectedResponse:   `{"exchangedAmount":"100000.99","exchangedRate":"1"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			path:               "/alive",
			method:             http.MethodGet,
			expectedResponse:   `{"status":true}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			path:               "/api/v1/currency-format",
			method:             http.MethodGet,
			expectedResponse:   `{"formats":[]}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			path:               "/api/v1/available-currencies",
			method:             http.MethodGet,
			expectedResponse:   `{"currencies":[]}`,
			expectedStatusCode: http.StatusOK,
		},
	} {
		t.Run(tt.path, func(t *testing.T) {
			logger := log.NewTest()
			mux := router.NewRouter(logger, mockService{}, mockService{})

			rr := httptest.NewRecorder()
			req := newRequest(t, tt.method, tt.path, tt.body)

			mux.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func newRequest(t *testing.T, method, path, bodyString string) *http.Request {
	t.Helper()

	var body io.Reader
	if bodyString != "" {
		body = bytes.NewReader([]byte(bodyString))
	}

	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	req.Header.Set("X-Correlation-ID", uuid.NewString())

	return req
}

type mockService struct{}

func (m mockService) GetForAmount(
	_ context.Context,
	_ decimal.Decimal,
	_ string,
	_ string,
) (decimal.Decimal, error) {
	return decimal.NewFromFloat(1.3523), nil
}

func (m mockService) UpsertExchangeRate(_ context.Context, _ domain.ExchangeRates) error {
	return nil
}

func (m mockService) ListExchangeRates(_ context.Context) ([]domain.ExchangeRates, error) {
	return []domain.ExchangeRates{}, nil
}

func (m mockService) Exchange(_ context.Context, request domain.ExchangeRequest) (*domain.Exchange, error) {
	return &domain.Exchange{
		Amount: request.Amount,
		Rate:   decimal.New(1, 0),
	}, nil
}

func (m mockService) GetFormats(_ context.Context) ([]formatdomain.CurrencyFormat, error) {
	return []formatdomain.CurrencyFormat{}, nil
}

func (m mockService) GetFormatted(_ context.Context, _ []formatdomain.FormatRequest) ([]formatdomain.FormatResponse, error) {
	return nil, nil
}

func (m mockService) GetByCode(_ context.Context, _ string) ([]formatdomain.CurrencyFormat, error) {
	return []formatdomain.CurrencyFormat{}, nil
}

func (m mockService) GetAvailableCurrenciesFromRates(_ context.Context) ([]string, error) {
	return []string{}, nil
}

func (m mockService) ExchangeMultiple(_ context.Context, _ domain.MultipleExchangesReq) (domain.MultipleExchangeResponse, error) {
	return domain.MultipleExchangeResponse{}, nil
}
*/
