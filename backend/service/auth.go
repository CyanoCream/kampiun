package service

import (
	"context"
	"errors"
	"time"

	"kampiun/domain"
	authjwt "kampiun/kernel/auth"
	"kampiun/kernel/notify"
	"kampiun/kernel/security"
)

type AuthService struct {
	repo     domain.Repository
	hasher   security.Hasher
	tokens   *authjwt.TokenMaker
	notifier notify.Notifier
}

func NewAuth(repo domain.Repository, hasher security.Hasher, tokens *authjwt.TokenMaker, notifier notify.Notifier) *AuthService {
	return &AuthService{repo: repo, hasher: hasher, tokens: tokens, notifier: notifier}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (string, domain.User, error) {
	in.Email = trimLower(in.Email)
	if in.Email == "" || in.Name == "" || len(in.Password) < 6 {
		return "", domain.User{}, domain.ErrInvalid
	}
	if _, err := s.repo.UserFindByEmail(in.Email); err == nil {
		return "", domain.User{}, domain.ErrDuplicate
	}
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return "", domain.User{}, err
	}
	u := domain.User{
		ID: domain.NewID(), Email: in.Email, Name: in.Name, Phone: in.Phone,
		PasswordHash: hash, Role: domain.RoleCustomer, CreatedAt: time.Now(),
	}
	if _, err := s.repo.UserCreate(u); err != nil {
		return "", domain.User{}, err
	}
	s.notifier.Notify(ctx, notify.Notification{Kind: notify.KindUserRegister,
		Title: "User baru daftar", Lines: []string{"Email: " + in.Email, "Nama: " + in.Name}})
	tok, err := s.tokens.Issue(u.ID, string(u.Role))
	return tok, u, err
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (string, domain.User, error) {
	u, err := s.repo.UserFindByEmail(in.Email)
	if err != nil {
		return "", domain.User{}, domain.ErrUnauthorized
	}
	if u.IsSuspended {
		return "", domain.User{}, errors.New("akun dinonaktifkan")
	}
	if !s.hasher.Verify(u.PasswordHash, in.Password) {
		return "", domain.User{}, domain.ErrUnauthorized
	}
	tok, err := s.tokens.Issue(u.ID, string(u.Role))
	return tok, u, err
}

func trimLower(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b = append(b, c)
	}
	return string(b)
}
