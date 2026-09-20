package pgdb

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kampiun/domain"
)

type Repo struct{ pool *pgxpool.Pool }

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

var _ domain.Repository = (*Repo)(nil)

const userCols = "id, email, name, phone, password_hash, role, is_suspended, created_at, updated_at"

// ---- User ----

func (r *Repo) UserCreate(u domain.User) (domain.User, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO users (id, email, name, phone, password_hash, role) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at, updated_at`,
		u.ID, strings.ToLower(u.Email), u.Name, u.Phone, u.PasswordHash, string(u.Role),
	).Scan(&u.CreatedAt, &u.UpdatedAt)
	return u, mapDup(err)
}

func (r *Repo) UserFindByEmail(email string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(context.Background(),
		`SELECT `+userCols+` FROM users WHERE email=$1 AND deleted_at IS NULL`,
		strings.ToLower(email),
	).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.PasswordHash, &u.Role, &u.IsSuspended, &u.CreatedAt, &u.UpdatedAt)
	return u, mapNotFound(err)
}

func (r *Repo) UserFindByID(id string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(context.Background(),
		`SELECT `+userCols+` FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.PasswordHash, &u.Role, &u.IsSuspended, &u.CreatedAt, &u.UpdatedAt)
	return u, mapNotFound(err)
}

func (r *Repo) UserUpdate(u domain.User) error {
	_, err := r.pool.Exec(context.Background(),
		`UPDATE users SET email=$2, name=$3, phone=$4, password_hash=$5, role=$6, is_suspended=$7, updated_at=now() WHERE id=$1`,
		u.ID, strings.ToLower(u.Email), u.Name, u.Phone, u.PasswordHash, string(u.Role), u.IsSuspended,
	)
	return err
}

// ---- Organization ----

func (r *Repo) OrgCreate(o domain.Organization) (domain.Organization, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO organizations (id, name, slug, owner_id) VALUES ($1,$2,$3,$4) RETURNING created_at, updated_at`,
		o.ID, o.Name, o.Slug, o.OwnerID,
	).Scan(&o.CreatedAt, &o.UpdatedAt)
	return o, mapDup(err)
}

func (r *Repo) OrgByID(id string) (domain.Organization, error) {
	var o domain.Organization
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM organizations WHERE id=$1`, id,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	return o, mapNotFound(err)
}

func (r *Repo) OrgBySlug(slug string) (domain.Organization, error) {
	var o domain.Organization
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM organizations WHERE slug=$1`, slug,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	return o, mapNotFound(err)
}

func (r *Repo) OrgListByOwner(ownerID string) ([]domain.Organization, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, name, slug, owner_id, created_at, updated_at FROM organizations WHERE owner_id=$1 ORDER BY created_at`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrgs(rows)
}

func scanOrgs(rows pgx.Rows) ([]domain.Organization, error) {
	var out []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.Slug, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// ---- Competition ----

const compCols = "id, org_id, name, slug, sport, format, is_public, access_code, status, created_at, updated_at"

const compScan = `SELECT ` + compCols + ` FROM competitions `

func (r *Repo) CompCreate(c domain.Competition) (domain.Competition, error) {
	fb, _ := sportFormatBytes(c.Format)
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO competitions (id, org_id, name, slug, sport, format, is_public, access_code, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at, updated_at`,
		c.ID, c.OrgID, c.Name, c.Slug, c.Sport, fb, c.IsPublic, c.AccessCode, string(c.Status),
	).Scan(&c.CreatedAt, &c.UpdatedAt)
	return c, mapDup(err)
}

func (r *Repo) CompByID(id string) (domain.Competition, error) {
	var c domain.Competition
	var fb []byte
	err := r.pool.QueryRow(context.Background(), compScan+`WHERE id=$1`, id).
		Scan(&c.ID, &c.OrgID, &c.Name, &c.Slug, &c.Sport, &fb, &c.IsPublic, &c.AccessCode, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	c.Format = sportFormatFromBytes(fb)
	return c, mapNotFound(err)
}

func (r *Repo) CompBySlug(slug string) (domain.Competition, error) {
	var c domain.Competition
	var fb []byte
	err := r.pool.QueryRow(context.Background(), compScan+`WHERE slug=$1`, slug).
		Scan(&c.ID, &c.OrgID, &c.Name, &c.Slug, &c.Sport, &fb, &c.IsPublic, &c.AccessCode, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	c.Format = sportFormatFromBytes(fb)
	return c, mapNotFound(err)
}

func (r *Repo) CompListByOrg(orgID string) ([]domain.Competition, error) {
	rows, err := r.pool.Query(context.Background(), compScan+`WHERE org_id=$1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComps(rows)
}

func (r *Repo) CompListPublic(ctx context.Context, sport, query string, limit, offset int) ([]domain.Competition, error) {
	q := `WHERE is_public`
	args := []any{}
	if sport != "" {
		args = append(args, sport)
		q += fmtAr(" AND sport=$%d", len(args))
	}
	if query != "" {
		args = append(args, "%"+query+"%")
		q += fmtAr(" AND name ILIKE $%d", len(args))
	}
	args = append(args, limit, offset)
	q += fmtAr(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args))
	rows, err := r.pool.Query(ctx, compScan+q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComps(rows)
}

func (r *Repo) CompUpdate(c domain.Competition) error {
	fb, _ := sportFormatBytes(c.Format)
	_, err := r.pool.Exec(context.Background(),
		`UPDATE competitions SET name=$2, sport=$3, format=$4, is_public=$5, access_code=$6, status=$7, updated_at=now() WHERE id=$1`,
		c.ID, c.Name, c.Sport, fb, c.IsPublic, c.AccessCode, string(c.Status),
	)
	return err
}

func scanComps(rows pgx.Rows) ([]domain.Competition, error) {
	var out []domain.Competition
	for rows.Next() {
		var c domain.Competition
		var fb []byte
		if err := rows.Scan(&c.ID, &c.OrgID, &c.Name, &c.Slug, &c.Sport, &fb, &c.IsPublic, &c.AccessCode, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Format = sportFormatFromBytes(fb)
		out = append(out, c)
	}
	return out, rows.Err()
}

// ---- Participant ----

func (r *Repo) PartCreate(p domain.Participant) (domain.Participant, error) {
	err := r.pool.QueryRow(context.Background(),
		`INSERT INTO participants (id, comp_id, name, seed, color, meta) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at`,
		p.ID, p.CompID, p.Name, p.Seed, p.Color, p.Meta,
	).Scan(&p.CreatedAt)
	return p, mapDup(err)
}

func (r *Repo) PartByID(id string) (domain.Participant, error) {
	var p domain.Participant
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, comp_id, name, seed, color, meta, created_at FROM participants WHERE id=$1`, id,
	).Scan(&p.ID, &p.CompID, &p.Name, &p.Seed, &p.Color, &p.Meta, &p.CreatedAt)
	return p, mapNotFound(err)
}

func (r *Repo) PartsByComp(compID string) ([]domain.Participant, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, comp_id, name, seed, color, meta, created_at FROM participants WHERE comp_id=$1 ORDER BY seed`, compID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Participant
	for rows.Next() {
		var p domain.Participant
		if err := rows.Scan(&p.ID, &p.CompID, &p.Name, &p.Seed, &p.Color, &p.Meta, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repo) PartUpdate(p domain.Participant) error {
	_, err := r.pool.Exec(context.Background(),
		`UPDATE participants SET name=$3, seed=$4, color=$5, meta=$6 WHERE id=$1 AND comp_id=$2`,
		p.ID, p.CompID, p.Name, p.Seed, p.Color, p.Meta,
	)
	return err
}

// ---- Match ----

func (r *Repo) MatchCreate(m domain.Match) (domain.Match, error) {
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO matches (id, comp_id, round, position, home_id, away_id, status, winner_id, is_bye) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		m.ID, m.CompID, m.Round, m.Position, m.HomeID, m.AwayID, string(m.Status), m.WinnerID, m.IsBye,
	)
	return m, err
}

func (r *Repo) MatchByID(id string) (domain.Match, error) {
	var m domain.Match
	err := r.pool.QueryRow(context.Background(),
		`SELECT id, comp_id, round, position, home_id, away_id, status, winner_id, started_at, finished_at, is_bye FROM matches WHERE id=$1`, id,
	).Scan(&m.ID, &m.CompID, &m.Round, &m.Position, &m.HomeID, &m.AwayID, &m.Status, &m.WinnerID, &m.StartedAt, &m.FinishedAt, &m.IsBye)
	return m, mapNotFound(err)
}

func (r *Repo) MatchesByComp(compID string) ([]domain.Match, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, comp_id, round, position, home_id, away_id, status, winner_id, started_at, finished_at, is_bye FROM matches WHERE comp_id=$1 ORDER BY round, position`, compID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Match
	for rows.Next() {
		var m domain.Match
		if err := rows.Scan(&m.ID, &m.CompID, &m.Round, &m.Position, &m.HomeID, &m.AwayID, &m.Status, &m.WinnerID, &m.StartedAt, &m.FinishedAt, &m.IsBye); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repo) MatchUpdate(m domain.Match) error {
	_, err := r.pool.Exec(context.Background(),
		`UPDATE matches SET round=$3, position=$4, home_id=$5, away_id=$6, status=$7, winner_id=$8, started_at=$9, finished_at=$10, is_bye=$11 WHERE id=$1 AND comp_id=$2`,
		m.ID, m.CompID, m.Round, m.Position, m.HomeID, m.AwayID, string(m.Status), m.WinnerID, m.StartedAt, m.FinishedAt, m.IsBye,
	)
	return err
}

// ---- Score events ----

func (r *Repo) EventAppend(e domain.ScoreEvent) error {
	var id string
	if e.ID == "" {
		id = newID()
	} else {
		id = e.ID
	}
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO score_events (id, match_id, side, delta, unit_index, by_user_id) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, e.MatchID, e.Side, e.Delta, e.UnitIndex, e.ByUserID,
	)
	return err
}

func (r *Repo) EventsByMatch(matchID string) ([]domain.ScoreEvent, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, match_id, side, delta, unit_index, by_user_id, created_at FROM score_events WHERE match_id=$1 ORDER BY created_at`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ScoreEvent
	for rows.Next() {
		var e domain.ScoreEvent
		if err := rows.Scan(&e.ID, &e.MatchID, &e.Side, &e.Delta, &e.UnitIndex, &e.ByUserID, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// ---- helpers ----

func newID() string { return domain.NewID() }

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func mapDup(err error) error {
	if err != nil && strings.Contains(err.Error(), "duplicate") {
		return domain.ErrDuplicate
	}
	return err
}
