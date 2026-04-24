package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/williamkoller/rest-api-sdd-go/config"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
	"github.com/williamkoller/rest-api-sdd-go/internal/domain/repository"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/cache"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountInactive    = errors.New("account inactive")
	ErrTokenInvalid       = errors.New("token invalid")
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *entity.User
}

type AuthUseCase struct {
	userRepo repository.UserRepository
	cache    cache.Cache
	cfg      *config.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, cache cache.Cache, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, cache: cache, cfg: cfg}
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*AuthTokens, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: login: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if !user.Active {
		return nil, ErrAccountInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return uc.issueTokens(ctx, user)
}

func (uc *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	// Check if token is invalidated in cache; miss errors mean "not revoked"
	cached, _ := uc.cache.Get(ctx, "refresh:revoked:"+refreshToken) //nolint:errcheck
	if cached != nil {
		return nil, ErrTokenInvalid
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(uc.cfg.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	userID, _ := claims["user_id"].(string)
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: refresh: %w", err)
	}
	if user == nil || !user.Active {
		return nil, ErrTokenInvalid
	}

	// Revoke old refresh token; cache failure is non-critical (token expires naturally)
	_ = uc.cache.Set(ctx, "refresh:revoked:"+refreshToken, []byte("1"), uc.cfg.JWT.RefreshTokenTTL) //nolint:errcheck

	return uc.issueTokens(ctx, user)
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	_ = uc.cache.Set(ctx, "refresh:revoked:"+refreshToken, []byte("1"), uc.cfg.JWT.RefreshTokenTTL) //nolint:errcheck
	return nil
}

func (uc *AuthUseCase) issueTokens(ctx context.Context, user *entity.User) (*AuthTokens, error) {
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"school_id": user.SchoolID,
		"role":      string(user.Role),
		"exp":       now.Add(uc.cfg.JWT.AccessTokenTTL).Unix(),
		"iat":       now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(uc.cfg.JWT.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("auth usecase: sign access token: %w", err)
	}

	refreshID := uuid.New().String()
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"jti":     refreshID,
		"exp":     now.Add(uc.cfg.JWT.RefreshTokenTTL).Unix(),
		"iat":     now.Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(uc.cfg.JWT.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("auth usecase: sign refresh token: %w", err)
	}

	_ = ctx
	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(uc.cfg.JWT.AccessTokenTTL.Seconds()),
		User:         user,
	}, nil
}
