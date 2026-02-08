package service

import (
	"context"
	"strings"
	"time"

	"auth-service/internal/model"
	"auth-service/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users  store.UserStore
	secret []byte
}

func NewAuthService(users store.UserStore, secret []byte) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterReq) (model.AuthResponse, error) {
	role := strings.ToUpper(strings.TrimSpace(req.Role))
	if role == "" {
		role = "TOURIST"
	}
	if role != "GUIDE" && role != "TOURIST" {
		return model.AuthResponse{}, ErrBadRole
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, ErrHashFailed
	}

	newID, err := s.users.NextUserID(ctx)
	if err != nil {
		return model.AuthResponse{}, ErrIDGenFailed
	}

	doc := model.UserDoc{
		ID:           newID,
		Email:        strings.TrimSpace(req.Email),
		Username:     strings.TrimSpace(req.Username),
		PasswordHash: string(pwHash),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := s.users.InsertUser(ctx, doc); err != nil {
		return model.AuthResponse{}, ErrDuplicateUser
	}

	token, err := s.issueToken(doc)
	if err != nil {
		return model.AuthResponse{}, ErrTokenFailed
	}

	return model.AuthResponse{ID: doc.ID, AccessToken: token}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginReq) (model.AuthResponse, error) {
	u, found, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil || !found {
		return model.AuthResponse{}, ErrInvalidCreds
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return model.AuthResponse{}, ErrInvalidCreds
	}

	token, err := s.issueToken(u)
	if err != nil {
		return model.AuthResponse{}, ErrTokenFailed
	}

	return model.AuthResponse{ID: u.ID, AccessToken: token}, nil
}

func (s *AuthService) GetByID(ctx context.Context, id int64) (model.UserDoc, bool, error) {
	return s.users.FindByID(ctx, id)
}

func (s *AuthService) issueToken(u model.UserDoc) (string, error) {
	claims := jwt.MapClaims{
		"id":               u.ID,
		"username":         u.Username,
		model.RoleClaimKey: u.Role,
		"exp":              time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}
