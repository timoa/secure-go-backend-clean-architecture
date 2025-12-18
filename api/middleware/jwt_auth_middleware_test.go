package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/middleware"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/tokenutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestJwtAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "access-secret"
	userObjectID := bson.NewObjectID()
	userID := userObjectID.Hex()

	accessToken, err := tokenutil.CreateAccessToken(&domain.User{ID: userObjectID, Name: "Test"}, secret, 1)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JwtAuthMiddleware(secret))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"userID": c.GetString("x-user-id")})
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var body map[string]string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Equal(t, userID, body["userID"])
	})

	t.Run("missing header", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JwtAuthMiddleware(secret))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var body domain.ErrorResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Equal(t, "Not authorized", body.Message)
	})

	t.Run("invalid token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JwtAuthMiddleware(secret))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("wrong secret", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JwtAuthMiddleware("wrong"))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

}
