package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"kampiun/domain"
	"kampiun/service"
)

type OrgAPI struct{ orgs *service.OrgService }

func NewOrgAPI(orgs *service.OrgService) *OrgAPI { return &OrgAPI{orgs: orgs} }

func (a *OrgAPI) HandleCreate(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "autentikasi diperlukan")
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	o, err := a.orgs.Create(r.Context(), service.CreateOrgInput{OwnerID: p.UserID, Name: in.Name})
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, o)
}

func (a *OrgAPI) HandleList(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "autentikasi diperlukan")
		return
	}
	orgs, err := a.orgs.ListByOwner(r.Context(), p.UserID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, orgs)
}

func statusOf(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrDuplicate):
		return http.StatusConflict
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
