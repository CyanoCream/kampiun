package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// migrate menjalankan skema Fase 1. Idempoten (CREATE TABLE IF NOT EXISTS).
func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, schema)
	return err
}

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id            text PRIMARY KEY,
    email         text NOT NULL UNIQUE,
    name          text NOT NULL,
    phone         text NOT NULL DEFAULT '',
    password_hash text NOT NULL,
    role          text NOT NULL DEFAULT 'customer',
    is_suspended  boolean NOT NULL DEFAULT false,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz
);

CREATE TABLE IF NOT EXISTS organizations (
    id         text PRIMARY KEY,
    name       text NOT NULL,
    slug       text NOT NULL UNIQUE,
    owner_id   text NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS competitions (
    id          text PRIMARY KEY,
    org_id      text NOT NULL REFERENCES organizations(id),
    name        text NOT NULL,
    slug        text NOT NULL UNIQUE,
    sport       text NOT NULL DEFAULT '',
    format      bytea,
    is_public   boolean NOT NULL DEFAULT false,
    access_code text NOT NULL DEFAULT '',
    status      text NOT NULL DEFAULT 'active',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS competitions_org_idx ON competitions(org_id);
CREATE INDEX IF NOT EXISTS competitions_public_idx ON competitions(is_public);

CREATE TABLE IF NOT EXISTS participants (
    id         text PRIMARY KEY,
    comp_id    text NOT NULL REFERENCES competitions(id),
    name       text NOT NULL,
    seed       int NOT NULL DEFAULT 0,
    color      text NOT NULL DEFAULT '#ccc',
    meta       text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS participants_comp_idx ON participants(comp_id);

CREATE TABLE IF NOT EXISTS matches (
    id          text PRIMARY KEY,
    comp_id     text NOT NULL REFERENCES competitions(id),
    round       int NOT NULL DEFAULT 1,
    position    int NOT NULL DEFAULT 0,
    home_id     text,
    away_id     text,
    status      text NOT NULL DEFAULT 'pending',
    winner_id   text,
    started_at  timestamptz,
    finished_at timestamptz,
    is_bye      boolean NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS matches_comp_idx ON matches(comp_id);

CREATE TABLE IF NOT EXISTS score_events (
    id         text PRIMARY KEY,
    match_id   text NOT NULL REFERENCES matches(id),
    side       text NOT NULL,
    delta      int NOT NULL DEFAULT 0,
    unit_index int NOT NULL DEFAULT 0,
    by_user_id text,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS events_match_idx ON score_events(match_id);
`
