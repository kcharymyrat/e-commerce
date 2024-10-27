package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func GenerateJWT(
	userID uuid.UUID, duration time.Duration, secretKey []byte, logger *zerolog.Logger,
) (string, *UserClaims, error) {
	user_claims, err := newUserClaims(userID, "", true, false, false, false, duration)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate token")
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user_claims)
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		logger.Error().Err(err).Msg("failed to sign token")
		return "", nil, err
	}

	return tokenStr, user_claims, nil
}

func ParseJWT(tokenString string, secretKey []byte, logger *zerolog.Logger) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		var msg string
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			msg = "malformed token"
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			msg = "invalid signature"
		case errors.Is(err, jwt.ErrTokenExpired):
			msg = "token expired"
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			msg = "token not valid yet"
		default:
			msg = "unknown error while parsing token"
		}
		logger.Error().Err(err).Msg(msg)
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		logger.Info().Str("user_id", claims.ID.String()).Msg("token is valid")
		return claims, nil
	}

	logger.Error().Msg("invalid token claims")
	return nil, errors.New("invalid token claims")
}
