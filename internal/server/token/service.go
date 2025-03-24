package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	secret     string
	expire     time.Duration
	signMethod jwt.SigningMethod
	tokenName  string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64 `json:"user_id,omitempty"`
}

// NewJwtService создает экземпляр jwt сервиса для авторизации
func NewJwtService(cfg *Config) *Service {
	return &Service{
		secret:     cfg.Secret,
		expire:     cfg.Expire,
		signMethod: jwt.SigningMethodHS256,
		tokenName:  cfg.Name,
	}
}

func (s *Service) LoadUserID(ctx context.Context) (domain.UserID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "unable to extract metadata")
	}

	values := md.Get(consts.AccessTokenHeader)
	if len(values) == 0 {
		return 0, status.Error(codes.Unauthenticated, "unable to extract authorization token")
	}

	token := values[0]

	claim, err := s.ParseToken(token)
	if err != nil {
		return 0, err
	}

	return domain.UserID(claim.UserID), nil
}

func (s *Service) BuildToken(id domain.UserID) (string, error) {
	token := jwt.NewWithClaims(s.signMethod, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expire)),
		},
		UserID: uint64(id),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ParseToken(tokenStr string) (*Claims, error) {
	cl := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, cl,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != s.signMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
			}

			return []byte(s.secret), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	return cl, nil
}
