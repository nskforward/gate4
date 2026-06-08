package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nskforward/gate4/internal/domain/model"
)

func Whoami(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(User).(model.User)
	if !ok {
		http.Error(w, "guest not allowed", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}
