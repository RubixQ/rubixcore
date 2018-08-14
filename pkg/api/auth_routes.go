package api

import (
	"net/http"
)

type authPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *App) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 	payload := new(authPayload)

		// 	err := json.NewDecoder(r.Body).Decode(payload)
		// 	if err != nil {
		// 		a.logger.Error("failed decoding authentication payload", zap.Error(err))
		// 		RenderBadRequest(w, err)
		// 		return
		// 	}

		// 	session := a.session.Copy()
		// 	defer session.Close()

		// 	repo := db.NewUserRepo(session)
		// 	user, err := repo.FindByCredentials(payload.Username, payload.Password)
		// 	if err != nil {
		// 		a.logger.Error("wrong username and/or password specified")
		// 		RenderBadRequest(w, fmt.Errorf("wrong username and/or password specified"))
		// 		return
		// 	}

		// 	now := time.Now()
		// 	claims := jwt.StandardClaims{
		// 		Issuer:    a.jwtIssuer,
		// 		IssuedAt:  now.Unix(),
		// 		ExpiresAt: now.Add(24 * 365 * time.Hour).Unix(),
		// 	}

		// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// 	tokenString, err := token.SignedString([]byte(a.jwtSecret))
		// 	if err != nil {
		// 		a.logger.Error("failed generating JWT after authentication")
		// 		RenderBadRequest(w, err)
		// 		return
		// 	}

		// 	RenderOk(w, struct {
		// 		Token string      `json:"token"`
		// 		User  interface{} `json:"user"`
		// 	}{
		// 		tokenString,
		// 		user,
		// 	})
	}
}
