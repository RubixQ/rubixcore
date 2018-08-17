package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/rubixq/rubixcore/pkg/api/repositories"
	"go.uber.org/zap"
)

type authPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticate verifies the credentials of a user and retusn
// a JWT for further authorization in subsequent calls to
// protected endpoints
func Authenticate(db *sqlx.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := new(authPayload)

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			logger.Error("failed decoding authentication payload", zap.Error(err))
			// TODO : Return bad request here
			return
		}

		repo := repositories.NewUserRepo(db)
		user, err := repo.FindByCredentials(payload.Username, payload.Password)
		if err != nil {
			logger.Error("failed finding system user by login credentials", zap.Error(err))
			// TODO : Return bad request here
			return
		}

		now := time.Now()
		claims := jwt.StandardClaims{
			Issuer:    "rubixcore",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("jwtsecret"))
		logger.Info("successfully generated jwt for user ", zap.Any("username", payload.Username), zap.Any("token", tokenStr))
		if err != nil {
			logger.Error("failed generating JWT after authentication", zap.Error(err))
			// TODO : Return bad request here
			return
		}

		_ = user
	}
}
