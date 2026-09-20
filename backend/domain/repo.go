package domain

import "context"

// Repository = port penyimpanan semua aggregate Fase 1 (satu implementasi Postgres).
type Repository interface {
	// User
	UserCreate(u User) (User, error)
	UserFindByEmail(email string) (User, error)
	UserFindByID(id string) (User, error)
	UserUpdate(u User) error

	// Organization
	OrgCreate(o Organization) (Organization, error)
	OrgByID(id string) (Organization, error)
	OrgBySlug(slug string) (Organization, error)
	OrgListByOwner(ownerID string) ([]Organization, error)

	// Competition
	CompCreate(c Competition) (Competition, error)
	CompByID(id string) (Competition, error)
	CompBySlug(slug string) (Competition, error)
	CompListByOrg(orgID string) ([]Competition, error)
	CompListPublic(ctx context.Context, sport, query string, limit, offset int) ([]Competition, error)
	CompUpdate(c Competition) error

	// Participant
	PartCreate(p Participant) (Participant, error)
	PartByID(id string) (Participant, error)
	PartsByComp(compID string) ([]Participant, error)
	PartUpdate(p Participant) error

	// Match
	MatchCreate(m Match) (Match, error)
	MatchByID(id string) (Match, error)
	MatchesByComp(compID string) ([]Match, error)
	MatchUpdate(m Match) error

	// Score events
	EventAppend(e ScoreEvent) error
	EventsByMatch(matchID string) ([]ScoreEvent, error)
}
