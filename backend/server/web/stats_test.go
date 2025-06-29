package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/internal/dto"
	"github.com/orion-tec/oriondns/internal/stats"
)

type MockStatsDB struct {
	mock.Mock
}

func (m *MockStatsDB) Insert(ctx context.Context, t time.Time, domain, domainType string) error {
	args := m.Called(ctx, t, domain, domainType)
	return args.Error(0)
}

func (m *MockStatsDB) GetMostUsedDomains(ctx context.Context, from, to time.Time, categories []string, limit int) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, categories, limit)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStatsDB) GetUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, domains []string) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, domains)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStatsDB) GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, categories []string) ([]stats.MostUsedDomainResponse, error) {
	args := m.Called(ctx, from, to, categories)
	return args.Get(0).([]stats.MostUsedDomainResponse), args.Error(1)
}

func (m *MockStatsDB) GetServerUsageByTimeRange(ctx context.Context, from, to time.Time, categories []string) ([]stats.ServerUsageByTimeRangeResponse, error) {
	args := m.Called(ctx, from, to, categories)
	return args.Get(0).([]stats.ServerUsageByTimeRangeResponse), args.Error(1)
}

func TestGetTimeFromFE(t *testing.T) {
	timestamp := int64(1640995200000)
	expectedTime := time.Unix(1640995200, 0).UTC()
	
	result := getTimeFromFE(timestamp)
	
	assert.Equal(t, expectedTime, result)
}

func TestHTTP_getMostUsedDomainsDashboard(t *testing.T) {
	mockStats := &MockStatsDB{}
	httpHandler := &HTTP{
		stats: mockStats,
	}

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	
	expectedResponse := []stats.MostUsedDomainResponse{
		{Domain: "google.com", Count: 100},
		{Domain: "facebook.com", Count: 50},
	}
	
	mockStats.On("GetMostUsedDomains", mock.Anything, from, to, []string{"search"}, 10).Return(expectedResponse, nil)

	reqBody := dto.GetMostUsedDomainsRequest{
		From:       from.Unix() * 1000,
		To:         to.Unix() * 1000,
		Categories: []string{"search"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/dashboard/most-used-domains", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	httpHandler.getMostUsedDomainsDashboard(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response []dto.GetMostUsedDomainsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response, 2)
	assert.Equal(t, "google.com", response[0].Domain)
	assert.Equal(t, int64(100), response[0].Count)
	assert.Equal(t, "facebook.com", response[1].Domain)
	assert.Equal(t, int64(50), response[1].Count)
	
	mockStats.AssertExpectations(t)
}

func TestHTTP_getMostUsedDomainsDashboard_EmptyCategories(t *testing.T) {
	mockStats := &MockStatsDB{}
	httpHandler := &HTTP{
		stats: mockStats,
	}

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	
	expectedResponse := []stats.MostUsedDomainResponse{
		{Domain: "google.com", Count: 100},
	}
	
	mockStats.On("GetMostUsedDomains", mock.Anything, from, to, []string{}, 10).Return(expectedResponse, nil)

	reqBody := dto.GetMostUsedDomainsRequest{
		From:       from.Unix() * 1000,
		To:         to.Unix() * 1000,
		Categories: []string{},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/dashboard/most-used-domains", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	httpHandler.getMostUsedDomainsDashboard(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	mockStats.AssertExpectations(t)
}

func TestHTTP_getServerUsageByTimeRangeDashboard(t *testing.T) {
	mockStats := &MockStatsDB{}
	httpHandler := &HTTP{
		stats: mockStats,
	}

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	
	expectedResponse := []stats.ServerUsageByTimeRangeResponse{
		{TimeRange: from, Count: 150},
		{TimeRange: from.Add(time.Hour), Count: 200},
	}
	
	mockStats.On("GetServerUsageByTimeRange", mock.Anything, from, to, []string{"search"}).Return(expectedResponse, nil)

	reqBody := dto.GetServerUsageByTimeRangeRequest{
		From:       from.Unix() * 1000,
		To:         to.Unix() * 1000,
		Categories: []string{"search"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/dashboard/server-usage-by-time-range", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	httpHandler.getServerUsageByTimeRangeDashboard(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response []dto.GetServerUsageByTimeRangeResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Len(t, response, 2)
	assert.Equal(t, int64(150), response[0].Count)
	assert.Equal(t, int64(200), response[1].Count)
	
	mockStats.AssertExpectations(t)
}

func TestHTTP_getMostUsedDomainsDashboard_InvalidJSON(t *testing.T) {
	mockStats := &MockStatsDB{}
	httpHandler := &HTTP{
		stats: mockStats,
	}

	req := httptest.NewRequest("POST", "/api/v1/dashboard/most-used-domains", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	httpHandler.getMostUsedDomainsDashboard(rr, req)
	
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHTTP_getServerUsageByTimeRangeDashboard_InvalidJSON(t *testing.T) {
	mockStats := &MockStatsDB{}
	httpHandler := &HTTP{
		stats: mockStats,
	}

	req := httptest.NewRequest("POST", "/api/v1/dashboard/server-usage-by-time-range", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	httpHandler.getServerUsageByTimeRangeDashboard(rr, req)
	
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}