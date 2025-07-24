package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/vistara-studio/vistara-be/internal/domain/user"
	"github.com/vistara-studio/vistara-be/pkg/cerr"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = cerr.New(fiber.StatusUnauthorized, "token expired", errors.New("invalid token"))
	ErrInvalidSignature = cerr.New(http.StatusUnauthorized, "different signature", errors.New("invalid signature"))
	ErrFailedClaimJWT   = cerr.New(http.StatusUnauthorized, "failed claim JWT", errors.New("invalid token"))
	ErrInvalidToken     = cerr.New(http.StatusUnauthorized, "invalid token", errors.New("invalid token"))
	ErrSignToken        = cerr.New(fiber.StatusInternalServerError, "failed to sign jwt-token", nil)
)

type Claims struct {
	UserID           string    `json:"user_id"`
	IsPremium        bool      `json:"is_premium"`
	PremiumExpiredAt time.Time `json:"premium_expired_at"`
	jwt.RegisteredClaims
}

type JWTStruct struct {
	secret []byte
}

func New(secret string) *JWTStruct {
	return &JWTStruct{
		secret: []byte(secret),
	}
}

func (j *JWTStruct) Encode(data *user.Table) (string, error) {
	claims := &Claims{
		UserID:           data.ID.String(),
		IsPremium:        data.IsPremium,
		PremiumExpiredAt: data.ExpiredAt,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "nusa",
			Subject:   "authentication",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", ErrSignToken.WithErr(err)
	}

	return signedToken, nil
}

func (j *JWTStruct) Decode(jwtToken string) (*Claims, error) {
	decoded, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return &Claims{}, ErrTokenExpired
		}

		return &Claims{}, ErrInvalidSignature
	}

	if !decoded.Valid {
		return &Claims{}, ErrInvalidToken
	}

	claims, ok := decoded.Claims.(*Claims)
	if !ok {
		return &Claims{}, ErrFailedClaimJWT
	}

	return claims, err
}
