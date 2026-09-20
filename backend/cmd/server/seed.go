package main

import (
	"context"

	"kampiun/domain"
	"kampiun/kernel/pgdb"
	"kampiun/kernel/security"
)

// seed membuat super admin awal bila belum ada.
func seed(ctx context.Context, repo *pgdb.Repo, hasher security.Hasher, email, password string) error {
	if _, err := repo.UserFindByEmail(email); err == nil {
		return nil
	}
	hash, err := hasher.Hash(password)
	if err != nil {
		return err
	}
	_, err = repo.UserCreate(domain.User{
		ID: domain.NewID(), Email: email, Name: "Super Admin",
		PasswordHash: hash, Role: domain.RoleSuperAdmin,
	})
	return err
}
