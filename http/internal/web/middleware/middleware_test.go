package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	mockuser "github.com/Lafetz/loyalty_marketplace/internal/web/mockUser"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/", RequireAuth())
	t.Run("Returns error if cookie is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		if err != nil {
			log.Fatal(err)
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

	})
	t.Run("Returns error if token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{
			Name:  "Authorization",
			Value: "token",
		}
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(cookie)

		if err != nil {
			log.Fatal(err)
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

	})
	t.Run("Successful if token is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		user := mockuser.User{
			Id:       uuid.New(),
			Username: "helloworld",
			Email:    "hellow@world.com",
		}
		token, err := jwtauth.CreateJwt(user)
		if err != nil {
			log.Fatal(err)
		}
		cookie := &http.Cookie{
			Name:  "Authorization",
			Value: token,
		}
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		req.AddCookie(cookie)

		if err != nil {
			log.Fatal(err)
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

	})

}
