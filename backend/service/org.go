package service

import (
	"context"
	"strings"

	"kampiun/domain"
)

// OrgService = kelola penyelenggara (organisasi).
type OrgService struct{ repo domain.Repository }

func NewOrgService(repo domain.Repository) *OrgService { return &OrgService{repo: repo} }

type CreateOrgInput struct {
	OwnerID string
	Name    string
}

func (s *OrgService) Create(ctx context.Context, in CreateOrgInput) (domain.Organization, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return domain.Organization{}, domain.ErrInvalid
	}
	o := domain.Organization{
		ID: domain.NewID(), Name: name, Slug: slugify(name), OwnerID: in.OwnerID,
	}
	return s.repo.OrgCreate(o)
}

func (s *OrgService) ListByOwner(ctx context.Context, ownerID string) ([]domain.Organization, error) {
	return s.repo.OrgListByOwner(ownerID)
}

func (s *OrgService) GetByID(ctx context.Context, id string) (domain.Organization, error) {
	return s.repo.OrgByID(id)
}

func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		out = domain.NewID()[:8]
	}
	return out
}
