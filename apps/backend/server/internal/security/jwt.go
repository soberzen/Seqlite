package security

import (
	"context"
	"errors"
	"fmt"
	"qingmo-server/internal/errs"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"

	JWTLeeway = 30 * time.Second

	keyRetentionMargin = time.Minute
)

type TokenManager struct {
	keyStore   SigningKeyStore
	issuer     string
	audience   string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Claims struct {
	TokenType TokenType `json:"token_type"`
	SessionID string    `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

func NewTokenManager(
	keyStore SigningKeyStore,
	issuer string,
	audience string,
	AccessTTL time.Duration,
	RefreshTTL time.Duration,
) *TokenManager {
	return &TokenManager{
		keyStore:   keyStore,
		issuer:     issuer,
		audience:   audience,
		AccessTTL:  AccessTTL,
		RefreshTTL: RefreshTTL,
	}
}

func (m *TokenManager) resolveTTL(tokenType TokenType) (time.Duration, error) {
	switch tokenType {
	case TokenTypeAccess:
		return m.AccessTTL, nil
	case TokenTypeRefresh:
		return m.RefreshTTL, nil
	default:
		return 0, errs.New(
			errs.KindInternal,
			fmt.Sprintf("unknown token type: %s", tokenType),
		)
	}
}

func (m *TokenManager) SignJWT(
	ctx context.Context,
	userID uint,
	tokenType TokenType,
	sessionID string,
) (string, *Claims, error) {
	if err := ctx.Err(); err != nil {
		return "", nil, err
	}

	ttl, err := m.resolveTTL(tokenType)

	if ttl <= 0 {
		return "", nil, errs.New(
			errs.KindInternal,
			fmt.Sprintf("invalid token TTL for %s", tokenType),
		)
	}

	now := time.Now()

	claims := &Claims{
		TokenType: tokenType,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings{m.audience},
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	retainUntil := claims.ExpiresAt.Time.
		Add(JWTLeeway).
		Add(keyRetentionMargin)

	kid, pair, err := m.keyStore.GetActiveKey(
		ctx,
		tokenType,
		retainUntil,
	)
	if err != nil {
		return "", nil, errs.Wrap(
			errs.KindInternal,
			"fail get active signing key",
			err,
		)
	}

	if kid == "" || pair == nil {
		return "", nil, errs.New(
			errs.KindInternal,
			"active signing key is missing",
		)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pair.PrivatePEM))
	if err != nil {
		return "", nil, errs.Wrap(
			errs.KindInternal,
			"fail parse RSA private key",
			err,
		)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil, errs.Wrap(
			errs.KindInternal,
			"fail sign JWT",
			err,
		)
	}

	return signed, claims, nil
}

func (m *TokenManager) VerifyJWT(
	ctx context.Context,
	token string,
	tokenType TokenType,
) (*Claims, error) {
	switch tokenType {
	case TokenTypeAccess:
	case TokenTypeRefresh:
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	claims := new(Claims)

	keyFunc := func(token *jwt.Token) (any, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errs.New(
				errs.KindUnauthorized,
				"token is missing kid header",
			)
		}

		pair, err := m.keyStore.GetKeyPair(ctx, tokenType, kid)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, errs.Wrap(
					errs.KindUnauthorized,
					"no active key for kid",
					err,
				)
			}
			return nil, errs.Wrap(errs.KindInternal, "fail get active key", err)
		}
		if pair == nil {
			return nil, errs.New(
				errs.KindInternal,
				fmt.Sprintf("verification key pair is missing for kid %s", kid),
			)
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pair.PublicPEM))
		if err != nil {
			return nil, errs.Wrap(
				errs.KindInternal,
				"fail parse RSA public key",
				err,
			)
		}
		return publicKey, nil
	}

	_, err := jwt.ParseWithClaims(
		token,
		claims,
		keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithNotBeforeRequired(),
		jwt.WithLeeway(JWTLeeway),
	)
	if err != nil {
		return nil, classifyJWTError(err)
	}
	if claims.TokenType != tokenType {
		return nil, errs.Wrap(
			errs.KindUnauthorized,
			fmt.Sprintf(
				"token type mismatch: got %s, want %s",
				claims.TokenType,
				tokenType,
			),
			err,
		)
	}

	return claims, nil
}

func classifyJWTError(err error) error {
	// 优先处理 ctx 取消，透传给上层特判
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return errs.Wrap(
			errs.KindUnauthorized,
			"token has expired",
			err,
		)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return errs.Wrap(errs.KindUnauthorized, "token is not valid yet", err)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return errs.Wrap(
			errs.KindUnauthorized,
			"token is malformed",
			err,
		)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return errs.Wrap(
			errs.KindUnauthorized,
			"token signature is invalid",
			err,
		)
	default:
		return err
	}
}
