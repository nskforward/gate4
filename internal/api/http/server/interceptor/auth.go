package interceptor

import (
	"context"
	"net/http"
	"strings"

	"github.com/nskforward/gate4/internal/api/http/server/handler"
	"github.com/nskforward/gate4/internal/domain/service"
)

type Auth struct {
	tokenService *service.TokenService
}

func NewAuth(tokenService *service.TokenService) *Auth {
	return &Auth{
		tokenService: tokenService,
	}
}

func (i *Auth) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := i.tokenService.FindUser(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), handler.User, user)))
	}
}
