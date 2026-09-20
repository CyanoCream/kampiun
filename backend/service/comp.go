package service

import (
	"context"
	"strings"

	"kampiun/domain"
)

// CompService = kelola kompetisi/turnamen & generate bracket.
type CompService struct{ repo domain.Repository }

func NewCompService(repo domain.Repository) *CompService { return &CompService{repo: repo} }

type CreateCompInput struct {
	OrgID  string         `json:"org_id"`
	Name   string         `json:"name"`
	Sport  string         `json:"sport"`
	Format domain.SportFormat `json:"format"`
	Public bool           `json:"public"`
}

func (s *CompService) Create(ctx context.Context, in CreateCompInput) (domain.Competition, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || in.OrgID == "" {
		return domain.Competition{}, domain.ErrInvalid
	}
	if in.Format.Sport == "" {
		// default set-based (voli/badminton/sepakbola): 3 set, sampai 25 point
		in.Format = domain.SportFormat{Sport: in.Sport, Unit: "set", TargetPoint: 25, BestOf: 3, UsesScore: true}
	}
	c := domain.Competition{
		ID: domain.NewID(), OrgID: in.OrgID, Name: name, Slug: slugify(name), Sport: in.Sport,
		Format: in.Format, IsPublic: in.Public, AccessCode: domain.ShortCode(4),
		Status: domain.CompActive,
	}
	return s.repo.CompCreate(c)
}

func (s *CompService) ListByOrg(ctx context.Context, orgID string) ([]domain.Competition, error) {
	return s.repo.CompListByOrg(orgID)
}

func (s *CompService) GetByID(ctx context.Context, id string) (domain.Competition, error) {
	return s.repo.CompByID(id)
}

func (s *CompService) ListPublic(ctx context.Context, sport, q string, limit, offset int) ([]domain.Competition, error) {
	return s.repo.CompListPublic(ctx, sport, q, limit, offset)
}

// CreateParticipants menambahkan N peserta & membangkitkan bracket knock-out.
func (s *CompService) CreateParticipants(ctx context.Context, compID string, names []string) ([]domain.Participant, error) {
	if len(names) < 2 {
		return nil, domain.ErrInvalid
	}
	parts := make([]domain.Participant, 0, len(names))
	for i, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			return nil, domain.ErrInvalid
		}
		p := domain.Participant{ID: domain.NewID(), CompID: compID, Name: n, Seed: i + 1, Color: defaultColor(i)}
		if _, err := s.repo.PartCreate(p); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	if err := s.generateBracket(ctx, compID, parts); err != nil {
		return nil, err
	}
	return parts, nil
}

// generateBracket membuat matches knock-out single-elimination + byes.
func (s *CompService) generateBracket(ctx context.Context, compID string, parts []domain.Participant) error {
	n := 1
	for n < len(parts) {
		n *= 2
	}
	// seeding: pasangan seed(i+1) vs seed(n-i)
	maxRound := 1
	x := n
	for x > 1 {
		x /= 2
		maxRound++
	}
	for round := 1; round <= maxRound; round++ {
		matches := n / (1 << round)
		for m := 0; m < matches; m++ {
			match := domain.Match{ID: domain.NewID(), CompID: compID, Round: round, Position: m,
				Status: domain.MatchPending}
			if round == 1 {
				hi := m*2
				lo := m*2 + 1
				if hi < len(parts) {
					match.HomeID = parts[hi].ID
				}
				if lo < len(parts) {
					match.AwayID = parts[lo].ID
				}
				if match.HomeID == "" {
					match.AwayID = "" // bye penuh
					match.IsBye = true
				}
			}
			if _, err := s.repo.MatchCreate(match); err != nil {
				return err
			}
		}
	}
	return nil
}

func defaultColor(i int) string {
	palette := []string{"#e74c3c", "#2ecc71", "#3498db", "#f39c12", "#9b59b6", "#1abc9c", "#e67e22", "#34495e"}
	return palette[i%len(palette)]
}
