package interceptor

import (
	"net/http"

	"github.com/nskforward/gate4/internal/domain/service"
)

type Auth struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewAuth(userService *service.UserService, tokenService *service.TokenService) *Auth {
	return &Auth{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (i *Auth) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
