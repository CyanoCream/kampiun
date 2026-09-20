package httpapi

import (
	"encoding/json"
	"net/http"

	"kampiun/domain"
	"kampiun/service"
)

type CompAPI struct {
	comps *service.CompService
	repo  domain.Repository
}

func NewCompAPI(comps *service.CompService, repo domain.Repository) *CompAPI {
	return &CompAPI{comps: comps, repo: repo}
}

func (a *CompAPI) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var in service.CreateCompInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	// pastikan pemilik punya org tsb
	if _, err := a.repo.OrgByID(in.OrgID); err != nil {
		writeErr(w, http.StatusNotFound, "organisasi tidak ditemukan")
		return
	}
	c, err := a.comps.Create(r.Context(), in)
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (a *CompAPI) HandleListByOrg(w http.ResponseWriter, r *http.Request) {
	comps, err := a.comps.ListByOrg(r.Context(), r.PathValue("org_id"))
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comps)
}

func (a *CompAPI) HandleAddParticipants(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Names []string `json:"names"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	parts, err := a.comps.CreateParticipants(r.Context(), r.PathValue("id"), in.Names)
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, parts)
}

// BracketNode = match dengan nama peserta (utk visualisasi bracket web/mobile).
type BracketNode struct {
	Match    domain.Match        `json:"match"`
	HomeName string             `json:"home_name"`
	AwayName string             `json:"away_name"`
}

// HandleBracket = bracket lengkap kompetisi (rounds → matches → nama peserta).
func (a *CompAPI) HandleBracket(w http.ResponseWriter, r *http.Request) {
	compID := r.PathValue("id")
	matches, err := a.repo.MatchesByComp(compID)
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	// nama peserta
	nameOf := map[string]string{}
	parts, _ := a.repo.PartsByComp(compID)
	for _, p := range parts {
		nameOf[p.ID] = p.Name
	}
	// kelompokkan per round
	byRound := map[int][]BracketNode{}
	var maxRound int
	for _, m := range matches {
		if m.Round > maxRound {
			maxRound = m.Round
		}
		byRound[m.Round] = append(byRound[m.Round], BracketNode{
			Match: m, HomeName: nameOf[m.HomeID], AwayName: nameOf[m.AwayID],
		})
	}
	rounds := make([]map[string]any, 0, maxRound)
	for r := 1; r <= maxRound; r++ {
		rounds = append(rounds, map[string]any{"round": r, "matches": byRound[r]})
	}
	writeJSON(w, http.StatusOK, rounds)
}

func (a *CompAPI) HandlePubicSearch(w http.ResponseWriter, r *http.Request) {
	comps, err := a.comps.ListPublic(r.Context(), r.URL.Query().Get("sport"), r.URL.Query().Get("q"), 50, 0)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comps)
}
