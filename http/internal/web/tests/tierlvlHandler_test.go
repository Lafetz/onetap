package handlers

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core "github.com/Lafetz/loyalty_marketplace/internal/loyalty"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	mockuser "github.com/Lafetz/loyalty_marketplace/internal/web/mockUser"
	"github.com/Lafetz/loyalty_marketplace/internal/web/routes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockNot struct {
}

func (m *mockNot) SendNotification(ctx context.Context, noti core.Notification) {}

func TestCreateTier(t *testing.T) {
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.InitAppRoutes(r, tierSvc, cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default()), slog.Default())

	merchantId := uuid.New()
	token, err := jwtauth.CreateJwt(mockuser.User{
		Id:       merchantId,
		Username: "uu",
		Email:    "hello@gmail.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name           string
		token          *http.Cookie
		requestBody    string
		expectedStatus int
	}{
		{
			name:  "Unauthorized - no token",
			token: nil,
			requestBody: `{
				"name": "Gold",
				"minpoint": 500
			}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Bad Request - invalid request body",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			requestBody:    `{ "name": "" }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Success - valid request",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			requestBody: `{
				"name": "yellow",
				"minpoint": 500
			}`,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/v1/merchants/"+merchantId.String()+"/tiers", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != nil {
				req.AddCookie(tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
func TestDeleteTier(t *testing.T) {
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.InitAppRoutes(r, tierSvc, cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default()), slog.Default())

	merchantId := uuid.New()
	token, err := jwtauth.CreateJwt(mockuser.User{
		Id:       merchantId,
		Username: "testuser",
		Email:    "testuser@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	err = tierSvc.CreateTierLevel(context.Background(), tier.NewTierLevel(merchantId, "Black", 5))
	if err != nil {

		log.Fatal(err)
	}
	tests := []struct {
		name           string
		token          *http.Cookie
		merchantId     string
		tierName       string
		expectedStatus int
	}{
		{
			name:           "Unauthorized - no token",
			token:          nil,
			merchantId:     merchantId.String(),
			tierName:       "Gold",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Bad Request - invalid merchant ID",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     "invalid-uuid",
			tierName:       "Gold",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Unauthorized - merchantId mismatch",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     uuid.New().String(),
			tierName:       "Gold",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Not Found - non-existent tier",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     merchantId.String(),
			tierName:       "NonExistentTier",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Success - valid request",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     merchantId.String(),
			tierName:       "Black",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/v1/merchants/"+tt.merchantId+"/tiers/"+tt.tierName, nil)

			if tt.token != nil {
				req.AddCookie(tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
func TestUpdateTierLevelHandler(t *testing.T) {
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.InitAppRoutes(r, tierSvc, cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default()), slog.Default())

	merchantId := uuid.New()
	token, err := jwtauth.CreateJwt(mockuser.User{
		Id:       merchantId,
		Username: "testuser",
		Email:    "testuser@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Create a test tier for the merchant
	testTierName := "Black"
	err = tierSvc.CreateTierLevel(context.Background(), tier.NewTierLevel(merchantId, testTierName, 5))
	if err != nil {
		t.Fatalf("Failed to create test tier: %v", err)
	}

	tests := []struct {
		name           string
		token          *http.Cookie
		body           string
		merchantId     uuid.UUID
		tierName       string
		expectedStatus int
	}{
		{
			name:           "Unauthorized - no token",
			token:          nil, // No token for unauthorized access
			body:           `{"name": "Gold", "minpoint": 500}`,
			merchantId:     merchantId,
			tierName:       testTierName,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Bad Request - missing fields",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			body:           `{"name": "Gold"}`, // Missing required fields
			merchantId:     merchantId,
			tierName:       testTierName,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Unauthorized - merchantId mismatch",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: uuid.New().String(),
			},
			body:           `{"name": "Gold", "minpoint": 500}`,
			merchantId:     uuid.New(), // Mismatched merchant ID
			tierName:       testTierName,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Not Found - tier does not exist",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			body:           `{"name": "NonExistentTier", "minpoint": 500}`,
			merchantId:     merchantId,
			tierName:       "NonExistentTier",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Success - valid request",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			body:           `{"name": "Gold", "minpoint": 500}`,
			merchantId:     merchantId,
			tierName:       testTierName,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/v1/merchants/"+merchantId.String()+"/tiers/"+tt.tierName, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			if tt.token != nil {
				req.AddCookie(tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
func TestGetTier(t *testing.T) {
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.InitAppRoutes(r, tierSvc, cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default()), slog.Default())

	merchantId := uuid.New()
	token, err := jwtauth.CreateJwt(mockuser.User{
		Id:       merchantId,
		Username: "testuser",
		Email:    "testuser@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Create a test tier for the merchant
	testTierName := "Pink"
	err = tierSvc.CreateTierLevel(context.Background(), tier.NewTierLevel(merchantId, testTierName, 5))
	if err != nil {
		t.Fatalf("Failed to create test tier: %v", err)
	}

	tests := []struct {
		name           string
		token          *http.Cookie
		merchantId     string
		tierName       string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{

		{
			name: "Bad Request - invalid merchant ID",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     "invalid-uuid",
			tierName:       testTierName,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"Error": "Invalid merchant ID",
			},
		},
		{
			name: "Not Found - tier does not exist",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     merchantId.String(),
			tierName:       "NonExistentTier",
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"Error": "Tier not found",
			},
		},
		{
			name: "Success - valid request",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     merchantId.String(),
			tierName:       testTierName,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"tier": map[string]interface{}{
					"name":     testTierName,
					"minpoint": 5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/merchants/"+tt.merchantId+"/tiers/"+tt.tierName, nil)

			if tt.token != nil {
				req.AddCookie(tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

		})
	}
}
func TestListTier(t *testing.T) {
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.InitAppRoutes(r, tierSvc, cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default()), slog.Default())

	merchantId := uuid.New()
	token, err := jwtauth.CreateJwt(mockuser.User{
		Id:       merchantId,
		Username: "testuser",
		Email:    "testuser@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Creating test tiers
	err = tierSvc.CreateTierLevel(context.Background(), tier.NewTierLevel(merchantId, "Gold", 500))
	if err != nil {
		log.Fatal(err)
	}
	err = tierSvc.CreateTierLevel(context.Background(), tier.NewTierLevel(merchantId, "Silver", 300))
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name           string
		token          *http.Cookie
		merchantId     string
		expectedStatus int
		expectedLength int
	}{

		{
			name: "Bad Request - invalid merchant ID",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedLength: 0,
		},
		{
			name: "Success - valid request",
			token: &http.Cookie{
				Name:  "Authorization",
				Value: token,
			},
			merchantId:     merchantId.String(),
			expectedStatus: http.StatusOK,
			expectedLength: 2, // Expected number of tiers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/merchants/"+tt.merchantId+"/tiers", nil)

			if tt.token != nil {
				req.AddCookie(tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Error decoding response: %v", err)
				}

				tiers, ok := response["tiers"].([]interface{})
				if !ok {
					t.Fatalf("Expected 'tiers' to be a list but got: %v", response["tiers"])
				}

				assert.Equal(t, tt.expectedLength, len(tiers))
			}
		})
	}
}
