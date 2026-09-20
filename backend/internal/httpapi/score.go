package httpapi

import (
	"encoding/json"
	"net/http"

	"kampiun/domain"
	"kampiun/service"
)

type ScoreAPI struct {
	repo  domain.Repository
	score *service.ScoreService
}

func NewScoreAPI(repo domain.Repository, score *service.ScoreService) *ScoreAPI {
	return &ScoreAPI{repo: repo, score: score}
}

// HandleScoreGet = live score publik (viewer), tanpa auth.
func (a *ScoreAPI) HandleScoreGet(w http.ResponseWriter, r *http.Request) {
	sc, err := a.score.Score(r.Context(), r.PathValue("id"))
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

// HandleTap = wasit menambah/mengurangi skor (auth required).
func (a *ScoreAPI) HandleTap(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "autentikasi diperlukan")
		return
	}
	var in service.ScoreInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "badan request tidak valid")
		return
	}
	in.MatchID = r.PathValue("id")
	in.ByUser = p.UserID
	sc, err := a.score.Apply(r.Context(), in)
	if err != nil {
		writeErr(w, statusOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sc)
}
