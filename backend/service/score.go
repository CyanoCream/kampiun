package service

import (
	"context"
	"errors"
	"time"

	"kampiun/domain"
)

// ScoreService = engine skor tap-tap inti. Stateless terhadap DB (hanya repo).
type ScoreService struct{ repo domain.Repository }

func NewScoreService(repo domain.Repository) *ScoreService { return &ScoreService{repo: repo} }

// Input membawa satu aksi skor pada match.
type ScoreInput struct {
	MatchID string
	Side    string `json:"side"`
	Delta   int    `json:"delta"`
	ByUser  string
}

// Apply menjalankan aksi skor, mengembalikan hasil (skor line) terbaru.
// Bila set/seri tuntas, menandai match finished & menimpa winner.
func (s *ScoreService) Apply(ctx context.Context, in ScoreInput) (MatchScore, error) {
	m, err := s.repo.MatchByID(in.MatchID)
	if err != nil {
		return MatchScore{}, err
	}
	if m.Status == domain.MatchFinished {
		return MatchScore{}, errors.New("pertandingan sudah selesai")
	}
	if in.Side != "home" && in.Side != "away" {
		return MatchScore{}, domain.ErrInvalid
	}
	c, err := s.repo.CompByID(m.CompID)
	if err != nil {
		return MatchScore{}, err
	}
	events, err := s.repo.EventsByMatch(m.ID)
	if err != nil {
		return MatchScore{}, err
	}
	ev := domain.ScoreEvent{
		ID: domain.NewID(), MatchID: m.ID, Side: in.Side, Delta: in.Delta, ByUserID: in.ByUser,
	}
	if err := s.repo.EventAppend(ev); err != nil {
		return MatchScore{}, err
	}
	events = append(events, ev)
	score := eval(c, m, events)
	if score.Finished {
		w := ""
		if score.HomeWon > score.AwayWon {
			w = m.HomeID
		} else {
			w = m.AwayID
		}
		now := time.Now()
		m.Status = domain.MatchFinished
		m.WinnerID = w
		m.FinishedAt = &now
		if err := s.repo.MatchUpdate(m); err != nil {
			return MatchScore{}, err
		}
	}
	return score, nil
}

// MatchScore = hasil evaluasi skor saat ini.
type MatchScore struct {
	HomeUnits, AwayUnits []int
	HomeWon, AwayWon     int
	Finished             bool
}

// eval menghitung skor dari events berdasarkan format dan kompetisi.
func eval(c domain.Competition, m domain.Match, events []domain.ScoreEvent) MatchScore {
	// hitung per sisi dari events; untuk mancing (UsesScore=false) nilai langsung dijumlah.
	homeUnits, awayUnits := []int{}, []int{}
	homeWon, awayWon := 0, 0
	curHome, curAway := 0, 0
	for _, e := range events {
		if e.Delta == 0 {
			continue
		}
		if c.Format.UsesScore {
			// set-based: +1 per tap
			if e.Side == "home" {
				curHome += e.Delta
			} else {
				curAway += e.Delta
			}
			// check finish set
			if c.Format.TargetPoint > 0 && (curHome >= c.Format.TargetPoint || curAway >= c.Format.TargetPoint) &&
				curHome != curAway {
				homeUnits = append(homeUnits, curHome)
				awayUnits = append(awayUnits, curAway)
				if curHome > curAway {
					homeWon++
				} else {
					awayWon++
				}
				curHome, curAway = 0, 0
			}
		} else {
			// mancing: +berat, tidak ada "set"
			if e.Side == "home" {
				curHome += e.Delta
			} else {
				curAway += e.Delta
			}
		}
	}
	if !c.Format.UsesScore {
		homeUnits = []int{curHome}
		awayUnits = []int{curAway}
	}
	finished := false
	if c.Format.UsesScore && c.Format.BestOf > 0 {
		finished = homeWon >= c.Format.BestOf/2+1 || awayWon >= c.Format.BestOf/2+1
	}
	return MatchScore{HomeUnits: homeUnits, AwayUnits: awayUnits, HomeWon: homeWon, AwayWon: awayWon, Finished: finished}
}

// Score returns current live score for a match (untuk viewer).
func (s *ScoreService) Score(ctx context.Context, matchID string) (MatchScore, error) {
	m, err := s.repo.MatchByID(matchID)
	if err != nil {
		return MatchScore{}, err
	}
	c, err := s.repo.CompByID(m.CompID)
	if err != nil {
		return MatchScore{}, err
	}
	events, err := s.repo.EventsByMatch(m.ID)
	if err != nil {
		return MatchScore{}, err
	}
	return eval(c, m, events), nil
}
