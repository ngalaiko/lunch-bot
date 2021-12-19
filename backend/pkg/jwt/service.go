package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lunch/pkg/jwt/keys"
	storage_keys "lunch/pkg/jwt/keys/storage"
	"lunch/pkg/users"

	"github.com/google/uuid"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Known errors.
var (
	ErrInvalidToken = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)

const (
	defaultIssuer = "lunch.bot"
	validFor      = 24 * time.Hour * 28 // 28 days
)

// Service allows to issue and verify jwt tokens.
type Service struct {
	keysDatabase storage_keys.Storage

	signer jose.Signer
}

// NewService creates a new jwt service.
func NewService(keysStorage storage_keys.Storage) *Service {
	return &Service{
		keysDatabase: keysStorage,
	}
}

type customClaims struct {
	Name string `json:"name"`
}

func (s *Service) init(ctx context.Context) error {
	if s.signer != nil {
		return fmt.Errorf("already initialized")
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	publicDER, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("failed to marshal encryption key: %w", err)
	}

	key, err := keys.New(publicDER)
	if err != nil {
		return fmt.Errorf("failed to create a key: %w", err)
	}

	if err := s.keysDatabase.Create(ctx, key); err != nil {
		return fmt.Errorf("failed to store key in the database: %w", err)
	}

	options := (&jose.SignerOptions{}).
		WithHeader("kid", key.ID).
		WithType("JWT")

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       privateKey,
	}, options)
	if err != nil {
		return err
	}

	s.signer = signer

	return nil
}

// NewToken creates a new signed JWT.
func (s *Service) NewToken(ctx context.Context, user *users.User) (*Token, error) {
	if s.signer == nil {
		if err := s.init(ctx); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	claims := &jwt.Claims{
		ID:       uuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  user.ID,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(validFor)),
	}
	customClaims := &customClaims{
		Name: user.Name,
	}

	token, err := jwt.Signed(s.signer).Claims(claims).Claims(customClaims).CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("failed to create a signed token: %w", err)
	}

	return &Token{
		Token:     token,
		User:      user,
		ExpiresAt: claims.Expiry.Time(),
	}, nil
}

// Verify checks token signature and returns it's meaningful content.
func (s *Service) Verify(ctx context.Context, token string) (*Token, error) {
	jwtoken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if len(jwtoken.Headers) == 0 {
		return nil, ErrInvalidToken
	}

	id := jwtoken.Headers[0].KeyID

	pubicKey, err := s.get(ctx, id)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	claims := &jwt.Claims{}
	customClaims := &customClaims{}
	if err := jwtoken.Claims(pubicKey, claims, customClaims); err != nil {
		return nil, ErrInvalidToken
	}

	validateErr := claims.ValidateWithLeeway(jwt.Expected{
		Time:   time.Now(),
		Issuer: defaultIssuer,
	}, time.Second)
	switch {
	case validateErr == nil:
		return &Token{
			Token: token,
			User: &users.User{
				ID:   claims.Subject,
				Name: customClaims.Name,
			},
			ExpiresAt: claims.Expiry.Time(),
		}, nil
	case errors.Is(validateErr, jwt.ErrExpired):
		return nil, ErrTokenExpired
	default:
		return nil, ErrInvalidToken
	}
}

func (s *Service) get(ctx context.Context, id string) (*ecdsa.PublicKey, error) {
	key, err := s.keysDatabase.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	untypedResult, err := x509.ParsePKIXPublicKey(key.PublicDER)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PKIX public key: %w", err)
	}

	switch v := untypedResult.(type) {
	case *ecdsa.PublicKey:
		return v, nil
	default:
		return nil, fmt.Errorf("unknown public key type: %T", v)
	}
}
