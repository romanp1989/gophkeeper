package interceptors

import (
	"context"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// authContext извлекает userID из JWT токена и добавляет его в контекст запроса
func authContext(tokenService *token.Service, ctx context.Context) (context.Context, error) {
	userID, err := tokenService.LoadUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user in claims")
	}

	ctx = context.WithValue(ctx, consts.UserIDKeyCtx, userID)

	return ctx, nil
}

// Authentication создает и возвращает interceptor для серверных вызовов gRPC.
// Автоматически применяется ко всем вызовам, кроме методов регистрации и входа в систему.
func Authentication(tokenService *token.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if strings.Contains(info.FullMethod, "Register") || strings.Contains(info.FullMethod, "Login") {
			return handler(ctx, req)
		}

		ctx, err := authContext(tokenService, ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
