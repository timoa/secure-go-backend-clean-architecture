package tokenutil_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/tokenutil"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCreateAccessTokenAndExtractIDFromToken(t *testing.T) {
	user := &domain.User{ID: bson.NewObjectID(), Name: "Test"}
	secret := "secret"

	token, err := tokenutil.CreateAccessToken(user, secret, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, err := tokenutil.ExtractIDFromToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), id)
}

func TestCreateRefreshTokenAndExtractIDFromToken(t *testing.T) {
	user := &domain.User{ID: bson.NewObjectID(), Name: "Test"}
	secret := "refresh-secret"

	token, err := tokenutil.CreateRefreshToken(user, secret, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, err := tokenutil.ExtractIDFromToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, user.ID.Hex(), id)
}

func TestIsAuthorized(t *testing.T) {
	user := &domain.User{ID: bson.NewObjectID(), Name: "Test"}
	secret := "secret"

	accessToken, err := tokenutil.CreateAccessToken(user, secret, 1)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		authorized, err := tokenutil.IsAuthorized(accessToken, secret)
		assert.NoError(t, err)
		assert.True(t, authorized)
	})

	t.Run("wrong secret", func(t *testing.T) {
		authorized, err := tokenutil.IsAuthorized(accessToken, "wrong")
		assert.Error(t, err)
		assert.False(t, authorized)
	})

	t.Run("invalid token", func(t *testing.T) {
		authorized, err := tokenutil.IsAuthorized("invalid", secret)
		assert.Error(t, err)
		assert.False(t, authorized)
	})
}

func TestExtractIDFromToken(t *testing.T) {
	user := &domain.User{ID: bson.NewObjectID(), Name: "Test"}
	secret := "secret"

	accessToken, err := tokenutil.CreateAccessToken(user, secret, 1)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		id, err := tokenutil.ExtractIDFromToken(accessToken, secret)
		assert.NoError(t, err)
		assert.Equal(t, user.ID.Hex(), id)
	})

	t.Run("wrong secret", func(t *testing.T) {
		id, err := tokenutil.ExtractIDFromToken(accessToken, "wrong")
		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("unexpected signing method", func(t *testing.T) {
		header, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT"})
		payload, _ := json.Marshal(map[string]string{"id": user.ID.Hex()})

		unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString([]byte("sig"))

		id, err := tokenutil.ExtractIDFromToken(unsigned, secret)
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("missing id claim", func(t *testing.T) {
		claims := jwt.MapClaims{"name": "no-id"}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(secret))
		assert.NoError(t, err)

		id, err := tokenutil.ExtractIDFromToken(signed, secret)
		assert.Error(t, err)
		assert.Empty(t, id)
	})
}
