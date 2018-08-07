package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
)

func (a *App) requireJWTAuthentication(nextHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authData := strings.Split(authHeader, " ")
		if len(authData) != 2 {
			a.logger.Error("authorization header is invalid")
			RenderBadRequest(w, fmt.Errorf("authorization header is invalid"))
			return
		}

		jwtToken := authData[1]
		claims, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing algorithm")
			}

			return []byte(a.jwtSecret), nil
		})

		if err != nil {
			a.logger.Error("failed validating jwt token", zap.Error(err))
			RenderBadRequest(w, fmt.Errorf("failed validating jwt token"))
			return
		}

		data := claims.Claims.(*jwt.StandardClaims)
		a.logger.Info("jwt claims", zap.Any("data", data))

		if data.Issuer != a.jwtIssuer {
			a.logger.Error("jwt is from the wrong issuer")
			RenderBadRequest(w, fmt.Errorf("jwt is from the wrong issuer"))
			return
		}

		if time.Now().After(time.Unix(data.ExpiresAt, 0)) {
			a.logger.Error("jwt has expired")
			RenderBadRequest(w, fmt.Errorf("jwt has expired"))
			return
		}

		nextHandler.ServeHTTP(w, r)
	}
}
