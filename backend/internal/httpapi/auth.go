package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"kampiun/domain"
	"kampiun/service"
)

type AuthAPI struct{ auth *service.AuthService }

func NewAuthAPI(auth *service.AuthService) *AuthAPI { return &AuthAPI{auth: auth} }

func (a *AuthAPI) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var in service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	tok, u, err := a.auth.Register(r.Context(), in)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, domain.ErrDuplicate):
			status = http.StatusConflict
		}
		writeErr(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"token": tok, "user": u})
}

func (a *AuthAPI) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var in service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	tok, u, err := a.auth.Login(r.Context(), in)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "email atau password salah")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": tok, "user": u})
}
